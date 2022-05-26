package timer

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
)

type queue struct {
	*tview.Table
	km    map[string]*help.Binding
	title string
}

// newQueue returns a new queue.
func newQueue() *queue {
	return &queue{
		Table:  tview.NewTable(),
		title: " Queue ",
		km: map[string]*help.Binding{
			"Select": help.NewBinding(
				help.WithKey(tcell.KeyEnter),
				help.WithHelp("Select"),
			),
		},
	}
}
