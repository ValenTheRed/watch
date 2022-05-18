package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Wtc struct {
	app  *tview.Application
	main *tview.TextView
}

func NewWtc(app *tview.Application) *Wtc {
	main := tview.NewTextView()
	main.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	return &Wtc{
		app:  app,
		main: main,
	}
}
