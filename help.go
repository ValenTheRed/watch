package main

import (
	"strings"

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

func (hv *HelpView) UpdateDisplay() {
	sep := " • "

	view := strings.Builder{}
	for _, bindings := range [][]*Binding{
		hv.locals,
		hv.globals,
	} {
		for _, b := range bindings {
			if !b.IsEnabled() {
				continue
			}
			var key string
			if b.Key() == tcell.KeyRune {
				key = string(b.Rune())
			} else {
				key = tcell.KeyNames[b.Key()]
			}
			view.WriteString(key + sep + b.Help())
		}
	}

	hv.SetText(view.String())
	view.Reset()
}
