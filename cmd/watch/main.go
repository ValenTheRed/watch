package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ValenTheRed/watch/internal/widget"
	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

var (
	usage = `usage: watch [-help] [duration]
A clock with a stopwatch and a timer.

Specify a duration to start a timer. Or, leave it alone to start a stopwatch.

optional arguments:
duration    supported formats - [[hh:]mm:]ss
-help	    display this help message and exit`
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}

	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight

	if err := clipboard.Init(); err != nil {
		panic(err)
	}
}

var (
	// SetTheme will be called before the running to application and
	// is expected to set the style of the widgets that will be
	// displayed.
	SetTheme func()

	ColorBackground = tcell.GetColor(colorful.Hcl(308.3, 0.02548, 0.04965).Hex())
	ColorForeground = tcell.GetColor(colorful.Hcl(0, 0.0001262, 0.8941).Hex())
	ColorPrimary    = tcell.GetColor(colorful.Hcl(15, .7, .5).Hex())
	ColorSecondary  = tcell.GetColor(colorful.Hcl(300, .5, .5).Hex())
	ColorBorder     = tcell.GetColor(colorful.Hcl(0, 4.714e-05, 0.2262).Hex())
	ColorSurface    = tcell.GetColor(colorful.Hcl(0, 6.055e-05, 0.336).Hex())
	ColorShadow     = tcell.ColorGrey
)

func main() {
	flag.Parse()

	durations := make([]int, len(flag.Args()))
	for i := range durations {
		var err error
		durations[i], err = ParseDuration(flag.Arg(i))
		if err != nil {
			log.Fatalln(fmt.Errorf("main: %v", err))
		}
		if durations[i] == 0 {
			log.Fatalln(fmt.Errorf("main: 0 not allowed; only positive integers"))
		}
	}

	app := tview.NewApplication().EnableMouse(true)
	if len(durations) == 0 {
		app = Stopwatch(app)
	} else {
		app = Timer(app, durations)
	}

	SetTheme()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// Stopwatch returns app after setting the root and starting the
// stopwatch.
func Stopwatch(app *tview.Application) *tview.Application {
	s := widget.NewStopwatch()
	s.Changed = func() {
		app.Draw()
	}
	l := widget.NewLapTable()

	type info struct {
		km     widget.KeyMap
		button *tview.Button
		action func()
	}
	interactions := struct {
		lap, playpause, restart, quit, copy info
	}{
		lap: info{
			km:     widget.KeyMap{Key: "l", Desc: "lap"},
			button: tview.NewButton("⚑ lap"),
		},
		restart: info{
			km:     widget.KeyMap{Key: "r", Desc: "restart"},
			button: tview.NewButton("● restart"),
		},
		playpause: info{
			km:     widget.KeyMap{Key: "space", Desc: "play/pause"},
			button: tview.NewButton("❚❚ pause"),
		},
		quit: info{
			km:     widget.KeyMap{Key: "q", Desc: "quit"},
			button: nil,
		},
		copy: info{
			km:     widget.KeyMap{Key: "y/c", Desc: "copy laps"},
			button: tview.NewButton(":: copy laps"),
		},
	}

	interactions.copy.action = func() {
		var lines []byte
		for row := l.GetRowCount() - 1; row > -1; row-- {
			lap, time, overall := l.GetLap(row)
			lines = append(lines, []byte(fmt.Sprintf("%2d", lap))...)
			lines = append(lines, ' ')
			lines = append(lines, []byte(widget.SecondWithColons(time))...)
			lines = append(lines, ' ')
			lines = append(lines, []byte(widget.SecondWithColons(overall))...)
			lines = append(lines, '\n')
		}
		clipboard.Write(clipboard.FmtText, lines)
	}
	interactions.lap.action = func() {
		l.AddLap(s.ElapsedSeconds())
	}
	interactions.restart.action = func() {
		s.Restart()
	}
	interactions.playpause.action = func() {
		if s.Running() {
			s.Stop()
		} else {
			s.Start()
		}
	}
	interactions.quit.action = func() {
		app.Stop()
	}

	var setSelectedButton = func(interaction info) {
		interaction.button.SetSelectedFunc(func() {
			interaction.action()
			// This simulates a button press.
			//
			// How?
			// Things to know:
			// - app.Draw() is called automatically after a Input/MouseHandler
			// - tview.Button's Draw() will use the highlight colors only
			// when it is in focus,
			// - tview.Button's mousehandler sets the focus to itself,
			// calls selected func in the same goroutine and then
			// returns.
			//
			// So, when we click a button, tview.Button's MouseHandler
			// get's called. This sets the focus to itself. As soon as
			// MouseHandler ends, the button updates and looks
			// highlighted. After __ milliseconds, the focus goes back
			// to another widget and the screen redraws. This
			// unhiglights the button.
			//
			// In a goroutine prevent deadlock.
			go func() {
				<-time.After(80 * time.Millisecond)
				app.SetFocus(l)
				app.Draw()
			}()
		})
	}
	setSelectedButton(interactions.lap)
	setSelectedButton(interactions.copy)
	setSelectedButton(interactions.restart)
	setSelectedButton(interactions.playpause)

	// Match playpause label to the action that the button will take
	// when pressed. Why not change labels inside of action() func of
	// playpause? Because, in the case when stopwatch is restarted, the
	// button label would not update.
	go func() {
		const (
			play  = "▶ play"
			pause = "❚❚ pause"
		)
		b := interactions.playpause.button
		for {
			if label := b.GetLabel(); s.Running() && label != pause {
				b.SetLabel(pause)
			} else if !s.Running() && label != play {
				b.SetLabel(play)
			}
		}
	}()

	bc := widget.NewButtonColumn([]*tview.Button{
		interactions.lap.button, interactions.playpause.button,
		interactions.restart.button, interactions.copy.button,
	})

	hv := widget.NewHelpView([]widget.KeyMap{
		interactions.lap.km, interactions.playpause.km,
		interactions.restart.km, interactions.quit.km, interactions.copy.km,
	})
	hv.SetDynamicColors(true)
	hv.SetTextAlign(tview.AlignCenter)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'l':
				interactions.lap.action()
				return nil
			case 'r':
				interactions.restart.action()
				return nil
			case 'q':
				interactions.quit.action()
				return nil
			case ' ':
				interactions.playpause.action()
				return nil
			case 'y', 'c':
				interactions.copy.action()
				return nil
			}
		}
		return event
	})

	f := tview.NewFlex().SetDirection(tview.FlexRow)
	f.AddItem(s, 0, 1, false)
	f.AddItem(bc, 0, 1, false)
	f.AddItem(hv, 2, 1, false)

	s.SetVerticalAlign(widget.AlignDown)
	s.SetBorderPadding(1, 1, 2, 2)
	bc.SetVerticalAlign(widget.AlignUp)
	bc.SetBorderPadding(1, 1, 2, 2)

	root := tview.NewFlex()
	root.AddItem(l, 0, 1, true)
	root.AddItem(f, 0, 3, false)

	SetTheme = func() {
		l.SetBorder(true)
		l.SetBorderColor(ColorSecondary)
		l.SetSelectedStyle(tcell.StyleDefault.Background(ColorPrimary))
		l.SetBackgroundColor(ColorBackground)
		l.SetHeaderStyle(tcell.StyleDefault.Foreground(ColorForeground))
		l.SetUnderlineStyle(tcell.StyleDefault.Foreground(ColorSecondary))
		l.SetCellStyle(tcell.StyleDefault.Foreground(ColorForeground))
		s.SetBackgroundColor(ColorBackground)
		bc.SetBackgroundColor(ColorBackground)
		hv.SetBackgroundColor(ColorBackground)
		s.TextColor = ColorForeground
		s.ShadowColor = ColorShadow
		hv.SetKeyStyle(tcell.StyleDefault.Foreground(ColorSurface))
		hv.SetDescStyle(tcell.StyleDefault.Foreground(ColorBorder))
		hv.SetSeparatorStyle(tcell.StyleDefault.Foreground(ColorBorder))
		var setButtonColor = func(b *tview.Button) {
			b.SetBackgroundColor(ColorPrimary)
			b.SetBackgroundColorActivated(ColorForeground)
			b.SetLabelColor(ColorForeground)
			b.SetLabelColorActivated(ColorPrimary)
		}
		setButtonColor(interactions.lap.button)
		setButtonColor(interactions.restart.button)
		setButtonColor(interactions.playpause.button)
		setButtonColor(interactions.copy.button)
	}

	s.Start()
	return app.SetRoot(root, true)
}

// Timer returns app after setting the root and starting the timer.
func Timer(app *tview.Application, durations []int) *tview.Application {
	t := widget.NewTimer(durations[0])
	p := widget.NewProgressBar()

	t.Changed = func() {
		p.SetPercent(t.ElapsedSeconds() * 100 / t.TotalSeconds())
		app.Draw()
	}

	q := widget.NewQueue(durations...)
	q.SetSelectedFunc(func(row int) {
		duration := q.GetCell(row, 1).GetReference().(int)
		t.SetTotalDuration(duration)
		t.Restart()
	})
	t.SetDoneFunc(func() {
		q.Next()
	})

	type info struct {
		km     widget.KeyMap
		button *tview.Button
		action func()
	}
	interactions := struct {
		prev, next, playpause, restart, quit info
	}{
		prev: info{
			km:     widget.KeyMap{Key: "p", Desc: "prev"},
			button: tview.NewButton("← prev"),
		},
		next: info{
			km:     widget.KeyMap{Key: "n", Desc: "next"},
			button: tview.NewButton("→ next"),
		},
		restart: info{
			km:     widget.KeyMap{Key: "r", Desc: "restart"},
			button: tview.NewButton("● restart"),
		},
		playpause: info{
			km:     widget.KeyMap{Key: "space", Desc: "play/pause"},
			button: tview.NewButton("❚❚ pause"),
		},
		quit: info{
			km:     widget.KeyMap{Key: "q", Desc: "quit"},
			button: nil,
		},
	}

	interactions.next.action = func() {
		q.Next()
	}
	interactions.prev.action = func() {
		q.Previous()
	}
	interactions.restart.action = func() {
		t.Restart()
	}
	interactions.playpause.action = func() {
		if t.Running() {
			t.Stop()
		} else {
			t.Start()
		}
	}
	interactions.quit.action = func() {
		app.Stop()
	}

	var setSelectedButton = func(interaction info) {
		interaction.button.SetSelectedFunc(func() {
			interaction.action()
			// This simulates a button press.
			//
			// How?
			// Things to know:
			// - app.Draw() is called automatically after a Input/MouseHandler
			// - tview.Button's Draw() will use the highlight colors only
			// when it is in focus,
			// - tview.Button's mousehandler sets the focus to itself,
			// calls selected func in the same goroutine and then
			// returns.
			//
			// So, when we click a button, tview.Button's MouseHandler
			// get's called. This sets the focus to itself. As soon as
			// MouseHandler ends, the button updates and looks
			// highlighted. After __ milliseconds, the focus goes back
			// to another widget and the screen redraws. This
			// unhiglights the button.
			//
			// In a goroutine prevent deadlock.
			go func() {
				<-time.After(80 * time.Millisecond)
				app.SetFocus(q)
				app.Draw()
			}()
		})
	}
	setSelectedButton(interactions.next)
	setSelectedButton(interactions.prev)
	setSelectedButton(interactions.restart)
	setSelectedButton(interactions.playpause)

	// Match playpause label to the action that the button will take
	// when pressed. Why not change labels inside of action() func of
	// playpause? Because, in the case when stopwatch is restarted, the
	// button label would not update.
	go func() {
		const (
			play  = "▶ play"
			pause = "❚❚ pause"
		)
		b := interactions.playpause.button
		for {
			if label := b.GetLabel(); t.Running() && label != pause {
				b.SetLabel(pause)
			} else if !t.Running() && label != play {
				b.SetLabel(play)
			}
		}
	}()

	bc := widget.NewButtonColumn([]*tview.Button{
		interactions.prev.button, interactions.playpause.button,
		interactions.restart.button, interactions.next.button,
	})

	hv := widget.NewHelpView([]widget.KeyMap{
		interactions.prev.km, interactions.playpause.km,
		interactions.restart.km, interactions.next.km, interactions.quit.km,
	})
	hv.SetDynamicColors(true)
	hv.SetTextAlign(tview.AlignCenter)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'p':
				interactions.prev.action()
				return nil
			case 'n':
				interactions.next.action()
				return nil
			case 'r':
				interactions.restart.action()
				return nil
			case 'q':
				interactions.quit.action()
				return nil
			case ' ':
				interactions.playpause.action()
				return nil
			}
		}
		return event
	})

	f := tview.NewFlex().SetDirection(tview.FlexRow)
	f.AddItem(t, 0, 2, false)
	f.AddItem(p, 0, 1, false)
	f.AddItem(bc, 0, 2, false)
	f.AddItem(hv, 2, 1, false)

	t.SetVerticalAlign(widget.AlignDown)
	t.SetBorderPadding(1, 1, 2, 2)
	p.SetAlign(widget.AlignCenter)
	p.SetBorderPadding(0, 0, 2, 2)
	bc.SetVerticalAlign(widget.AlignUp)
	bc.SetBorderPadding(1, 1, 2, 2)

	root := tview.NewFlex()
	root.AddItem(q, 0, 1, true)
	root.AddItem(f, 0, 3, false)

	SetTheme = func() {
		q.SetBorder(true)
		q.SetBorderColor(ColorSecondary)
		q.SetBackgroundColor(ColorBackground)
		q.SetSelectedStyle(tcell.StyleDefault.Background(ColorPrimary))
		q.SetHeaderStyle(tcell.StyleDefault.Foreground(ColorForeground))
		q.SetUnderlineStyle(tcell.StyleDefault.Foreground(ColorSecondary))
		q.SetCellStyle(tcell.StyleDefault.Foreground(ColorForeground))
		t.SetBackgroundColor(ColorBackground)
		p.SetBackgroundColor(ColorBackground)
		bc.SetBackgroundColor(ColorBackground)
		hv.SetBackgroundColor(ColorBackground)
		t.TextColor = ColorForeground
		t.ShadowColor = ColorShadow
		p.TextColor = ColorForeground
		p.ShadowColor = ColorShadow
		hv.SetKeyStyle(tcell.StyleDefault.Foreground(ColorSurface))
		hv.SetDescStyle(tcell.StyleDefault.Foreground(ColorBorder))
		hv.SetSeparatorStyle(tcell.StyleDefault.Foreground(ColorBorder))
		var setButtonColor = func(b *tview.Button) {
			b.SetBackgroundColor(ColorPrimary)
			b.SetBackgroundColorActivated(ColorForeground)
			b.SetLabelColor(ColorForeground)
			b.SetLabelColorActivated(ColorPrimary)
		}
		setButtonColor(interactions.prev.button)
		setButtonColor(interactions.next.button)
		setButtonColor(interactions.restart.button)
		setButtonColor(interactions.playpause.button)
	}

	t.Start()
	return app.SetRoot(root, true)
}

// ParseDuration returns the total number of seconds in dur, which must
// be of format [[hh:]mm:]ss.
func ParseDuration(dur string) (int, error) {
	var hr, min, sec int

	if m, err := regexp.MatchString(`^\d*$`, dur); m {
		if err != nil {
			return 0, err
		}
		sec, _ = strconv.Atoi(dur)
	} else if m, err := regexp.MatchString(`^\d+:\d{2}$`, dur); m {
		if err != nil {
			return 0, err
		}
		s := strings.Split(dur, ":")
		min, _ = strconv.Atoi(s[0])
		sec, _ = strconv.Atoi(s[1])
		// it's okay for minute field to be more than 60
		if err = checkField(sec, 0); err != nil {
			return 0, err
		}
	} else if m, err := regexp.MatchString(`^\d+:\d{2}:\d{2}$`, dur); m {
		if err != nil {
			return 0, err
		}
		s := strings.Split(dur, ":")
		hr, _ = strconv.Atoi(s[0])
		min, _ = strconv.Atoi(s[1])
		sec, _ = strconv.Atoi(s[2])
		if err = checkField(sec, min); err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("duration must be in [[hh:]mm:]ss format")
	}

	return (hr * 3600) + (min * 60) + sec, nil
}

// checkField returns error if sec/min field are not less than 60.
func checkField(sec, min int) error {
	var errmsg string
	if sec >= 60 {
		errmsg = "second's"
	}
	if min >= 60 && errmsg == "" {
		errmsg = "minute's"
	} else if min >= 60 {
		errmsg += " and minute's"
	}
	if errmsg != "" {
		return fmt.Errorf("%v field must be less than 60", errmsg)
	}
	return nil
}
