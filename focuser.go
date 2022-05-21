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

		// FIXME: UpdateDisplay is bugged. Or rather, my approach to the
		// way all the widgets update the their text is flawed.
		//
		// If `wtc.help` executes this function then, the main
		// application goroutine gets blocked and the program halts. Why
		// does this happen? Speculating: probably because of the way
		// TextView widgets -- which `wtc.help` is -- run the handler
		// installed using `SetChangedFunc`.
		// [Concurrency in tview](https://github.com/rivo/tview/wiki/Concurrency)
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

		// FIXME: UpdateDisplay is bugged. Same as `focusFunc()`
		wtc.help.UpdateDisplay()
	}
}
