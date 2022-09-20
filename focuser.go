package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
)

type Focuser interface {
	Title() string
	SetTitle(string) *tview.Box

	SetBorderColor(tcell.Color) *tview.Box
	SetTitleColor(tcell.Color) *tview.Box
}

func focusFunc(widget Focuser, km help.KeyMaper) func() {
	return func() {
		widget.
			SetTitle("[" + widget.Title() + "]").
			SetTitleColor(tcell.NewHexColor(0xf06db7)).
			SetBorderColor(tcell.NewHexColor(0xf06db7))
		watch.help.SetLocals(km)
		watch.help.UpdateDisplay()
	}
}

func blurFunc(widget Focuser) func() {
	return func() {
		widget.
			SetTitle(widget.Title()).
			SetTitleColor(tview.Styles.TitleColor).
			SetBorderColor(tview.Styles.BorderColor)
		watch.help.UnsetLocals()
		watch.help.UpdateDisplay()
	}
}
