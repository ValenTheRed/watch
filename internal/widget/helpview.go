package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeyMap struct {
	Key, Desc string
}

type HelpView struct {
	*tview.TextView

	// keys is a keymap.
	keys []KeyMap

	// Styles for the a key's shortcut, it's help description, and the
	// separator character between two keymap.
	keyStyle, descStyle, separatorStyle tcell.Style
}
