package laps

import (
	"bytes"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/stopwatch"
	"github.com/ValenTheRed/watch/utils"
)

type keyMap struct {
	Lap, Yank, Copy *help.Binding
}

type Laps struct {
	*tview.Table
	km    keyMap
	title string

	sw  *stopwatch.Stopwatch
	app *tview.Application
}

func New(sw *stopwatch.Stopwatch, app *tview.Application) *Laps {
	err := clipboard.Init()
	if err != nil {
		// panicing since subsequent calls to clipboard functions are
		// going to panic anyway.
		panic(err)
	}

	l := &Laps{
		title: " Lap ",
		Table: tview.NewTable(),

		app: app,
		sw:  sw,
		km: keyMap{
			Lap: help.NewBinding(
				help.WithRune('l'), help.WithHelp("Lap"),
			),
			Copy: help.NewBinding(
				help.WithRune('c'), help.WithHelp("Copy all"),
			),
			Yank: help.NewBinding(
				help.WithRune('y'), help.WithHelp("Copy row"),
			),
		},
	}

	l.
		initFirstRow().
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case l.km.Lap.Rune():
				// [Concurrency in tview](https://github.com/rivo/tview/wiki/Concurrency#event-handlers)
				// mentions not needing to call any redrawing functions
				// as application's main loop will do that for us.
				l.Lap()
			case l.km.Copy.Rune():
				l.Copy()
			case l.km.Yank.Rune():
				l.Yank()
			}
			return event
		}).
		SetTitle(l.title)

	return l
}

// initFirstRow inserts the first row of table with headings indicating:
// lap number, lap time and overall time elapsed when lapped.
func (l *Laps) initFirstRow() *Laps {
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

func (l *Laps) Title() string {
	return l.title
}

func (l *Laps) Keys() []*help.Binding {
	return []*help.Binding{l.km.Lap, l.km.Copy, l.km.Yank}
}

// Lap inserts new row in l. The row is inserted below the topmost i.e.
// header row.
func (l *Laps) Lap() {
	const row = 1

	var (
		overall    = l.sw.Elapsed()
		i, lapTime int
	)

	if l.GetRowCount() == row {
		i, lapTime = 1, overall
	} else {
		i = l.GetCell(row, 0).Reference.(int) + 1
		lapTime = overall - l.GetCell(row, 2).Reference.(int)
	}

	l.InsertRow(row)
	l.SetCell(row, 0, newLapCell(fmt.Sprintf("%02d", i), i))
	l.SetCell(row, 1, newLapCell(utils.FormatSecond(lapTime), lapTime))
	l.SetCell(row, 2, newLapCell(utils.FormatSecond(overall), overall))
}

// newLapCell returns a Lap cell with common style applied.
func newLapCell(text string, ref interface{}) *tview.TableCell {
	return tview.NewTableCell(text).
		SetReference(ref).
		SetAlign(tview.AlignCenter)
}

// Copy copies all of the rows into the system clipboard, except for the
// first row.
func (l *Laps) Copy() {
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
func (l *Laps) Yank() {
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
