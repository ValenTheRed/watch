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
			SetTitleColor(tcell.ColorOrange).
			SetBorderColor(tcell.ColorOrange)
		wtc.help.SetLocals(km)
		wtc.help.UpdateDisplay()
	}
}

func blurFunc(widget Focuser) func() {
	return func() {
		widget.
			SetTitle(widget.Title()).
			SetTitleColor(tview.Styles.TitleColor).
			SetBorderColor(tview.Styles.BorderColor)
		wtc.help.UnsetLocals()
		wtc.help.UpdateDisplay()
	}
}
