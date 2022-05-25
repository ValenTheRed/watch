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

// New returns a new Stopwatch.
func New(app *tview.Application) *Stopwatch {
	return &Stopwatch{
		app:  app,
		Flex: tview.NewFlex(),
		swtc: newStopwatch(),
		laps: newLaps(),
	}
}

// Init initialises Stopwatch and it's components. Should be run
// immediately after New().
func (sw *Stopwatch) Init() *Stopwatch {
	sw.swtc.init()
	sw.laps.init()

	sw.swtc.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case sw.swtc.km["Reset"].Rune():
			sw.ResetStopwatch()
			sw.swtc.km["Start"].SetDisable(true)
			sw.swtc.km["Stop"].SetDisable(false)
		case sw.swtc.km["Stop"].Rune():
			if sw.swtc.km["Stop"].IsEnabled() {
				sw.Stop()
				sw.swtc.km["Stop"].SetDisable(true)
				sw.swtc.km["Start"].SetDisable(false)
			}
		case sw.swtc.km["Start"].Rune():
			if sw.swtc.km["Start"].IsEnabled() {
				sw.Start()
				sw.swtc.km["Start"].SetDisable(true)
				sw.swtc.km["Stop"].SetDisable(false)
			}
		}
		return event
	})

	sw.laps.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case sw.laps.km["Lap"].Rune():
			// [Concurrency in tview](https://github.com/rivo/tview/wiki/Concurrency#event-handlers)
			// mentions not needing to call any redrawing functions
			// as application's main loop will do that for us.
			// TODO: sw.Lap()
		case sw.laps.km["Copy"].Rune():
			sw.laps.copy()
		case sw.laps.km["Yank"].Rune():
			sw.laps.yank()
		case sw.laps.km["Reset"].Rune():
			sw.laps.Clear()
			sw.laps.initFirstRow()
			// The selection stays at the row after a clear. So, the
			// first row is unselected after a clear. Row selection
			// appears back again only when:
			// - more rows have been added, the previous row is
			// selected, or
			// - standard Table movement keys are used
			// So, we reselect the first row.
			sw.laps.Select(0, 0)
		}
		return event
	})

	// TODO: focus capture

	sw.SetDirection(tview.FlexRow).
		// Fix size
		AddItem(sw.swtc, 3, 0, true).
		// Flexible size
		AddItem(sw.laps, 0, 1, false)

	sw.QueueStopwatchDraw()
	return sw
}

// QueueStopwatchDraw queues sw's stopwatch component for redraw.
func (sw *Stopwatch) QueueStopwatchDraw() {
	go sw.app.QueueUpdateDraw(func() {
		sw.swtc.SetText(utils.FormatSecond(sw.swtc.elapsed))
	})
}

func (sw *Stopwatch) Start() {
	if !sw.swtc.running {
		sw.swtc.running = true
		go utils.Worker(func() {
			sw.swtc.elapsed++
			sw.QueueStopwatchDraw()
		}, sw.swtc.stopMsg)
	}
}

func (sw *Stopwatch) Stop() {
	if sw.swtc.running {
		sw.swtc.running = false
		sw.swtc.stopMsg <- struct{}{}
	}
}

// ResetStopwatch reset the stopwatch back to 0 and starts it again. The
// stopwatch starts running even if it was paused prior to invocation.
func (sw *Stopwatch) ResetStopwatch() {
	sw.Stop()
	sw.swtc.elapsed = 0
	sw.QueueStopwatchDraw()
	sw.Start()
}
