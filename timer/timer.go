package timer

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

//go:embed "ping.flac"
var pingFile []byte

type keyMap struct {
	Reset, Stop, Start *help.Binding
}

// Timer component for Timer.
type timer struct {
	*tview.TextView
	km                map[string]*help.Binding
	title             string
	duration, elapsed int
	running           bool
	stopMsg           chan struct{}
}

// newTimer returns a new timer.
func newTimer(duration int) *timer {
	return &timer{
		TextView: tview.NewTextView(),
		title:    " Timer ",
		// Channel is buffered because: `Stop()` -- which sends on
		// `stopMsg` -- will be called by the instance of `worker()`
		// started by `Start()`, which has it's `quit` channel
		// set to `stopMsg`; `Stop()` will block an unbuffered `stopMsg`.
		stopMsg:  make(chan struct{}, 1),
		duration: duration,
		elapsed:  0,
		km: map[string]*help.Binding{
			"Reset": help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			"Pause": help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			"Start": help.NewBinding(
				help.WithRune('s'),
				help.WithHelp("Start"),
				help.WithDisable(true),
			),
		},
	}
}

// init returns an initialised t. Should be run immediately after
// newTimer().
func (t *timer) init() *timer {
	t.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(t.title)
	return t
}

// Title returns the default title of t.
func (t *timer) Title() string {
	return t.title
}

func (t *timer) String() string {
	const (
		boundary = "┃"
		fill     = "#"
	)

	elapsed := utils.FormatSecond(t.elapsed)
	dur := utils.FormatSecond(t.duration)

	_, _, width, _ := t.GetInnerRect()
	// +2 is for the boundary chars at the either end of the progress
	// bar.
	width -= len(elapsed) + 2*tview.TabSize + 2 + 2*tview.TabSize + len(dur)
	percent := t.elapsed * 100 / t.duration
	fillLen := width * percent / 100

	return fmt.Sprintf(
		"\t%s\t%s\t%s\t", elapsed,
		strings.Join([]string{
			boundary,
			strings.Repeat(fill, fillLen),
			strings.Repeat(" ", width-fillLen),
			boundary,
		}, ""),
		dur,
	)
}

func (t *timer) IsTimeLeft() bool {
	return t.elapsed < t.duration
}

type Timer struct {
	*tview.Flex
	app *tview.Application

	Timer *timer
}

// New returns a new Timer.
func New(duration int, app *tview.Application) *Timer {
	return &Timer{
		Flex:  tview.NewFlex(),
		app:   app,
		Timer: newTimer(duration),
	}
}

// Init initialises components of Timer. Must be run immediately after New().
func (t *Timer) Init() *Timer {
	t.Timer.init()

	t.Timer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case t.Timer.km["Reset"].Rune():
			t.Reset()
			t.Timer.km["Start"].SetDisable(true)
			t.Timer.km["Stop"].SetDisable(false)
		case t.Timer.km["Stop"].Rune():
			if t.Timer.km["Stop"].IsEnabled() {
				t.Stop()
				t.Timer.km["Stop"].SetDisable(true)
				t.Timer.km["Start"].SetDisable(false)
			}
		case t.Timer.km["Start"].Rune():
			if t.Timer.km["Start"].IsEnabled() {
				t.Start()
				t.Timer.km["Start"].SetDisable(true)
				t.Timer.km["Stop"].SetDisable(false)
			}
		}
		return event
	})

	t.SetDirection(tview.FlexRow).
		AddItem(t.Timer, 3, 0, true)

	t.UpdateDisplay()
	return t
}

// Keys returns the key bindings t uses for it's timer component.
func (t *Timer) Keys() []*help.Binding {
	return []*help.Binding{
		t.Timer.km["Start"],
		t.Timer.km["Stop"],
		t.Timer.km["Reset"],
	}
}

func (t *Timer) UpdateDisplay() {
	go t.app.QueueUpdateDraw(func() {
		t.Timer.SetText(t.Timer.String())
	})
}

func (t *Timer) Start() {
	if t.Timer.running || !t.Timer.IsTimeLeft() {
		return
	}
	t.Timer.running = true
	go utils.Worker(func() {
		t.Timer.elapsed++
		if !t.Timer.IsTimeLeft() {
			t.Stop()
			go Ping(bytes.NewReader(pingFile))
			t.Timer.km["Stop"].SetDisable(true)
		}
		t.UpdateDisplay()
	}, t.Timer.stopMsg)
}

func (t *Timer) Stop() {
	if t.Timer.running {
		t.Timer.running = false
		t.Timer.stopMsg <- struct{}{}
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.Timer.elapsed = 0
	t.UpdateDisplay()
	t.Start()
}

func Ping(r io.Reader) error {
	streamer, format, err := flac.Decode(r)
	if err != nil {
		return err
	}
	defer streamer.Close()

	done := make(chan struct{})
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- struct{}{}
	})))

	<-done
	return nil
}
