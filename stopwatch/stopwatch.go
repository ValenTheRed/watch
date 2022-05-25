package stopwatch

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

// Stopwatch widget component for Stopwatch.
type stopwatch struct {
	*tview.TextView
	km      map[string]*help.Binding
	title   string
	elapsed int
	running bool
	stopMsg chan struct{}
}

// newStopwatch returns a new stopwatch.
func newStopwatch() *stopwatch {
	return &stopwatch{
		TextView: tview.NewTextView(),
		stopMsg:  make(chan struct{}),
		title:    " Stopwatch ",
		km: map[string]*help.Binding{
			"Reset": help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			"Stop": help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			"Start": help.NewBinding(
				help.WithRune('s'), help.WithHelp("Start"),
				help.WithDisable(true),
			),
		},
	}
}

// init returns an initialised s. Should be run immediately after
// newStopwatch.
func (s *stopwatch) init() *stopwatch {
	s.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(s.title)
	return s
}

// Title returns the title of s.
func (s *stopwatch) Title() string {
	return s.title
}

// Keys returns the list of key bindings attached to s.
func (s *stopwatch) Keys() []*help.Binding {
	return []*help.Binding{
		s.km["Pause"], s.km["Start"], s.km["Reset"],
	}
}

type Stopwatch struct {
	*tview.Flex
	app *tview.Application

	swtc *stopwatch
	laps *laps
}

func New(app *tview.Application) *Stopwatch {
	sw := &Stopwatch{
		app:  app,
		Flex: tview.NewFlex(),
		swtc: newStopwatch(),
		laps: newLaps(),
	}

	sw.
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case sw.km.Reset.Rune():
				sw.Reset()
				sw.km.Start.SetDisable(true)
				sw.km.Stop.SetDisable(false)
			case sw.km.Stop.Rune():
				if sw.km.Stop.IsEnabled() {
					sw.Stop()
					sw.km.Stop.SetDisable(true)
					sw.km.Start.SetDisable(false)
				}
			case sw.km.Start.Rune():
				if sw.km.Start.IsEnabled() {
					sw.Start()
					sw.km.Start.SetDisable(true)
					sw.km.Stop.SetDisable(false)
				}
			}
			return event
		})

	sw.UpdateDisplay()
	return sw
}

func (sw *Stopwatch) Title() string {
	return sw.title
}

func (sw *Stopwatch) Elapsed() int {
	return sw.elapsed
}

func (sw *Stopwatch) Keys() []*help.Binding {
	return []*help.Binding{sw.km.Reset, sw.km.Start, sw.km.Stop}
}

func (sw *Stopwatch) UpdateDisplay() {
	go sw.app.QueueUpdateDraw(func() {
		sw.SetText(utils.FormatSecond(sw.elapsed))
	})
}

func (sw *Stopwatch) Start() {
	if !sw.running {
		sw.running = true
		go utils.Worker(func() {
			sw.elapsed++
			sw.UpdateDisplay()
		}, sw.stopMsg)
	}
}

func (sw *Stopwatch) Stop() {
	if sw.running {
		sw.running = false
		sw.stopMsg <- struct{}{}
	}
}

func (sw *Stopwatch) Reset() {
	sw.Stop()
	sw.elapsed = 0
	sw.UpdateDisplay()
	sw.Start()
}
