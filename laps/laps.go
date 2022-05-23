package laps

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/stopwatch"
	"github.com/ValenTheRed/watch/utils"
)

type keyMap struct {
	Lap, Copy *help.Binding
}

type Laps struct {
	*tview.Table
	km    keyMap
	title string

	sw  *stopwatch.Stopwatch
	app *tview.Application
}

func New(sw *stopwatch.Stopwatch, app *tview.Application) *Laps {
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
				// TODO: Copy to system clipboard
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
			SetExpansion(1),
	)
	l.SetCell(0, 1,
		newLapCell("Lap time", nil).
			SetExpansion(2),
	)
	l.SetCell(0, 2,
		newLapCell("Overall time", nil).
			SetExpansion(2),
	)
	return l
}

func (l *Laps) Title() string {
	return l.title
}

func (l *Laps) Keys() []*help.Binding {
	return []*help.Binding{l.km.Lap, l.km.Copy}
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
	l.SetCell(row, 0, newLapCell(fmt.Sprintf("%2d", i), i))
	l.SetCell(row, 1, newLapCell(utils.FormatSecond(lapTime), lapTime))
	l.SetCell(row, 2, newLapCell(utils.FormatSecond(overall), overall))
}

// newLapCell returns a Lap cell with common style applied.
func newLapCell(text string, ref interface{}) *tview.TableCell {
	return tview.NewTableCell(text).
		SetReference(ref).
		SetAlign(tview.AlignCenter)
}
