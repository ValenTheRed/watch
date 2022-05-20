package main

import (
	"github.com/rivo/tview"
)

type HelpView struct {
	*tview.TextView
	globals, locals []*Binding
	title           string
}
