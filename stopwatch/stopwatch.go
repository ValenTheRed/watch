package stopwatch

import (
	"fmt"

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
			"Pause": help.NewBinding(
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

	Swtc *stopwatch
	Laps *laps
}

// New returns a new Stopwatch.
func New(app *tview.Application) *Stopwatch {
	return &Stopwatch{
		app:  app,
		Flex: tview.NewFlex(),
		Swtc: newStopwatch(),
		Laps: newLaps(),
	}
}

// Init initialises Stopwatch and it's components. Should be run
// immediately after New().
func (sw *Stopwatch) Init() *Stopwatch {
	sw.Swtc.init()
	sw.Laps.init()

	sw.Swtc.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case sw.Swtc.km["Reset"].Rune():
			sw.ResetStopwatch()
			sw.Swtc.km["Start"].SetDisable(true)
			sw.Swtc.km["Pause"].SetDisable(false)
		case sw.Swtc.km["Pause"].Rune():
			if sw.Swtc.km["Pause"].IsEnabled() {
				sw.Stop()
				sw.Swtc.km["Pause"].SetDisable(true)
				sw.Swtc.km["Start"].SetDisable(false)
			}
		case sw.Swtc.km["Start"].Rune():
			if sw.Swtc.km["Start"].IsEnabled() {
				sw.Start()
				sw.Swtc.km["Start"].SetDisable(true)
				sw.Swtc.km["Pause"].SetDisable(false)
			}
		}
		return event
	})

	sw.Laps.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case sw.Laps.km["Lap"].Rune():
			// [Concurrency in tview](https://github.com/rivo/tview/wiki/Concurrency#event-handlers)
			// mentions not needing to call any redrawing functions
			// as application's main loop will do that for us.
			sw.Lap()
		case sw.Laps.km["Copy"].Rune():
			sw.Laps.copy()
		case sw.Laps.km["Yank"].Rune():
			sw.Laps.yank()
		case sw.Laps.km["Reset"].Rune():
			sw.Laps.reset()
		}
		return event
	})

	sw.SetDirection(tview.FlexRow).
		// Fix size
		AddItem(sw.Swtc, 3, 0, true).
		// Flexible size
		AddItem(sw.Laps, 0, 1, false)

	sw.QueueStopwatchDraw()
	return sw
}

// QueueStopwatchDraw queues sw's stopwatch component for redraw.
func (sw *Stopwatch) QueueStopwatchDraw() {
	go sw.app.QueueUpdateDraw(func() {
		sw.Swtc.SetText(utils.FormatSecond(sw.Swtc.elapsed))
	})
}

func (sw *Stopwatch) Start() {
	if !sw.Swtc.running {
		sw.Swtc.running = true
		go utils.Worker(func() {
			sw.Swtc.elapsed++
			sw.QueueStopwatchDraw()
		}, sw.Swtc.stopMsg)
	}
}

func (sw *Stopwatch) Stop() {
	if sw.Swtc.running {
		sw.Swtc.running = false
		sw.Swtc.stopMsg <- struct{}{}
	}
}

// ResetStopwatch reset the stopwatch back to 0 and starts it again. The
// stopwatch starts running even if it was paused prior to invocation.
func (sw *Stopwatch) ResetStopwatch() {
	sw.Stop()
	sw.Swtc.elapsed = 0
	sw.QueueStopwatchDraw()
	sw.Start()
}

// Lap creates a new lap entry.
func (sw *Stopwatch) Lap() {
	const row = 1

	var (
		overall    = sw.Swtc.elapsed
		i, lapTime int
	)

	if sw.Laps.GetRowCount() == row {
		i, lapTime = 1, overall
	} else {
		i = sw.Laps.GetCell(row, 0).Reference.(int) + 1
		lapTime = overall - sw.Laps.GetCell(row, 2).Reference.(int)
	}

	sw.Laps.InsertRow(row)
	sw.Laps.SetCell(row, 0, newLapCell(fmt.Sprintf("%02d", i), i))
	sw.Laps.SetCell(row, 1, newLapCell(utils.FormatSecond(lapTime), lapTime))
	sw.Laps.SetCell(row, 2, newLapCell(utils.FormatSecond(overall), overall))
}
