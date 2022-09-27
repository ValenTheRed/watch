package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ValenTheRed/watch/internal/widget"
	"github.com/rivo/tview"
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
}

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
		lap, playpause, restart, quit info
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
			km: widget.KeyMap{Key: "q", Desc: "quit"},
			button: nil,
		},
	}

	interactions.lap.action = func() {
		l.AddLap(s.ElapsedSeconds())
		if !l.HasFocus() {
			app.SetFocus(l)
		}
	}
	interactions.restart.action = func() {
		s.Restart()
		if !l.HasFocus() {
			app.SetFocus(l)
		}
	}
	interactions.playpause.action = func() {
		if s.Running() {
			s.Stop()
		} else {
			s.Start()
		}
		if !l.HasFocus() {
			app.SetFocus(l)
		}
	}
	interactions.quit.action = func() {
		app.Stop()
	}
	interactions.lap.button.SetSelectedFunc(interactions.lap.action)
	interactions.restart.button.SetSelectedFunc(interactions.restart.action)
	interactions.playpause.button.SetSelectedFunc(interactions.playpause.action)

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
		interactions.restart.button,
	})

	hv := widget.NewHelpView([]widget.KeyMap{
		interactions.lap.km, interactions.playpause.km,
		interactions.restart.km, interactions.quit.km,
	})
	hv.SetDynamicColors(true)
	hv.SetTextAlign(tview.AlignCenter)

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
			km: widget.KeyMap{Key: "q", Desc: "quit"},
			button: nil,
		},
	}

	interactions.next.action = func() {
		q.Next()
		if !q.HasFocus() {
			app.SetFocus(q)
		}
	}
	interactions.prev.action = func() {
		q.Previous()
		if !q.HasFocus() {
			app.SetFocus(q)
		}
	}
	interactions.restart.action = func() {
		t.Restart()
		if !q.HasFocus() {
			app.SetFocus(q)
		}
	}
	interactions.playpause.action = func() {
		if t.Running() {
			t.Stop()
		} else {
			t.Start()
		}
		if !q.HasFocus() {
			app.SetFocus(q)
		}
	}
	interactions.quit.action = func() {
		app.Stop()
	}
	interactions.next.button.SetSelectedFunc(interactions.next.action)
	interactions.prev.button.SetSelectedFunc(interactions.prev.action)
	interactions.restart.button.SetSelectedFunc(interactions.restart.action)
	interactions.playpause.button.SetSelectedFunc(interactions.playpause.action)

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
