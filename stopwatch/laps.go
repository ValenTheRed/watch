package stopwatch

import (
	"bytes"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"

	"github.com/ValenTheRed/watch/help"
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

// init returns an initialised l. Also initialises package clipboard.
// Should be run immediately after newLaps.
func (l *laps) init() *laps {
	l.
		initFirstRow().
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(l.title)

	err := clipboard.Init()
	if err != nil {
		// panicing since subsequent calls to clipboard functions are
		// going to panic anyway.
		panic(err)
	}

	return l
}

// Title returns the title of l.
func (l *laps) Title() string {
	return l.title
}

// Keys returns the list of key bindings attached to l.
func (l *laps) Keys() []*help.Binding {
	return []*help.Binding{
		l.km["Lap"], l.km["Copy"], l.km["Yank"], l.km["Reset"],
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

// Copy copies all of the rows into the system clipboard, except for the
// first row.
func (l *laps) copy() {
	lines := make([][]byte, 0, l.GetRowCount()-1)
	for row := 1; row < l.GetRowCount(); row++ {
		line := make([][]byte, 3, 3)
		for col := 0; col < 3; col++ {
			line[col] = []byte(l.GetCell(row, col).Text)
		}
		lines = append(lines, bytes.Join(line, []byte{byte(' ')}))
	}
	clipboard.Write(clipboard.FmtText, bytes.Join(lines, []byte{byte('\n')}))
}

// Yank() copies currently selected row, except if it is the first row.
func (l *laps) yank() {
	row, _ := l.GetSelection()
	// Return if header row
	if row == 0 {
		return
	}
	line := make([][]byte, 3, 3)
	for col := 0; col < 3; col++ {
		line[col] = []byte(l.GetCell(row, col).Text)
	}
	clipboard.Write(clipboard.FmtText, bytes.Join(line, []byte{byte(' ')}))
}

// reset clears l of all entries.
func (l *laps) reset() {
	l.Clear()
	l.initFirstRow()
	// The selection stays at the row after a clear. So, the
	// first row is unselected after a clear. Row selection
	// appears back again only when:
	// - more rows have been added, the previous row is
	// selected, or
	// - standard Table movement keys are used
	// So, we reselect the first row.
	l.Select(0, 0)
}
