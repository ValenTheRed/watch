package stopwatch

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

// Laps widget component for Stopwatch.
type laps struct {
	*tview.Table
	km    map[string]*help.Binding
	title string
}

// newLaps returns a new laps.
func newLaps() *laps {
	return &laps{
		Table: tview.NewTable(),
		title: " Laps ",
		km: map[string]*help.Binding{
			"Lap": help.NewBinding(
				help.WithRune('l'), help.WithHelp("Lap"),
			),
			"Copy": help.NewBinding(
				help.WithRune('c'), help.WithHelp("Copy all"),
			),
			"Yank": help.NewBinding(
				help.WithRune('y'), help.WithHelp("Copy row"),
			),
			"Reset": help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
		},
	}
}

// Keys returns the list of key bindings attached to l.
func (l *laps) Keys() []*help.Binding {
	return []*help.Binding{
		l.km["Lap"], l.km["Copy"], l.km["yank"], l.km["Reset"],
	}
}

// initFirstRow inserts and initialises the first row i.e. the row with
// column headers.
func (l *laps) initFirstRow() *laps {
	if l.GetRowCount() > 0 {
		return l
	}
	l.InsertRow(0)
	// SetExpansion applies on the whole column.
	l.SetCell(0, 0,
		newLapCell("Lap", nil).
			SetAttributes(tcell.AttrBold).
			SetExpansion(1),
	)
	l.SetCell(0, 1,
		newLapCell("Lap time", nil).
			SetAttributes(tcell.AttrBold).
			SetExpansion(2),
	)
	l.SetCell(0, 2,
		newLapCell("Overall time", nil).
			SetAttributes(tcell.AttrBold).
			SetExpansion(2),
	)
	return l
}

// newLapCell returns a Table cell with a default style for a laps cell
// applied.
func newLapCell(text string, ref interface{}) *tview.TableCell {
	return tview.NewTableCell(text).
		SetReference(ref).
		SetAlign(tview.AlignCenter)
}

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

// Keys returns the list of key bindings attached to s.
func (s *stopwatch) Keys() []*help.Binding {
	return []*help.Binding{
		s.km["Pause"], s.km["Start"], s.km["Reset"],
	}
}

type keyMap struct {
	Reset, Stop, Start *help.Binding
}

type Stopwatch struct {
	*tview.TextView
	elapsed int
	running bool
	stopMsg chan struct{}
	title   string
	km      keyMap

	app *tview.Application
}

func New(app *tview.Application) *Stopwatch {
	sw := &Stopwatch{
		app:      app,
		TextView: tview.NewTextView(),
		stopMsg:  make(chan struct{}),
		title:    " Stopwatch ",
		km: keyMap{
			Reset: help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			Stop: help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			Start: help.NewBinding(
				help.WithRune('s'),
				help.WithHelp("Start"),
				help.WithDisable(true),
			),
		},
	}

	sw.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
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
		}).
		SetTitle(sw.title)

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
