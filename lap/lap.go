package lap

import (
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
