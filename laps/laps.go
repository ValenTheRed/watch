package laps

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/stopwatch"
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
				// TODO: Add new row with lap info
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
	l.SetCell(0, 0, tview.NewTableCell("Lap"))
	l.SetCell(0, 1, tview.NewTableCell("Lap time"))
	l.SetCell(0, 2, tview.NewTableCell("Overall time"))
	return l
}

func (l *Laps) Title() string {
	return l.title
}

func (l *Laps) Keys() []*help.Binding {
	return []*help.Binding{l.km.Lap, l.km.Copy}
}
