package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HelpView struct {
	*tview.TextView
	globals, locals []*Binding
	title           string
}

func NewHelpView() *HelpView {
	hv := &HelpView{
		TextView: tview.NewTextView(),
		title:    " Help ",
	}
	hv.
		SetChangedFunc(func() {
			wtc.app.Draw()
		}).
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(focusFunc(hv)).
		SetBlurFunc(blurFunc(hv)).
		SetTitle(hv.title)

	return hv
}

func (hv *HelpView) Title() string {
	return hv.title
}
