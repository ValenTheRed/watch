package lap

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/stopwatch"
)

type keyMap struct {
	Lap, Copy *help.Binding
}

type Lap struct {
	*tview.Table
	km    keyMap
	title string

	sw  *stopwatch.Stopwatch
	app *tview.Application
}

func New(sw *stopwatch.Stopwatch, app *tview.Application) *Lap {
	l := &Lap{
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
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetTitleAlign(tview.AlignLeft).
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
