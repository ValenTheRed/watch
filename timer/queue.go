package timer

import (
	"github.com/ValenTheRed/watch/help"
	"github.com/rivo/tview"
)

type queue struct {
	*tview.Table
	km    map[string]*help.Binding
	title string
}
