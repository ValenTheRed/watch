package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Wtc struct {
	app  *tview.Application
	main *tview.TextView
	help *tview.TextView
	panels []*tview.TextView
}

func NewWtc(app *tview.Application) *Wtc {
	main := tview.NewTextView()
	main.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	help := tview.NewTextView()
	help.
		SetTextAlign(tview.AlignCenter).
		SetTitle("Help").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	return &Wtc{
		app:  app,
		main: main,
		help: help,
		panels: []*tview.TextView{ main, help },
	}
}

func (w *Wtc) Run() error {
	return w.app.SetRoot(w.setLayout(), true).Run()
}

func (w *Wtc) setLayout() *tview.Flex {
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(w.main, 0, 9, true).
		AddItem(w.help, 0, 1, false)
}

func (w *Wtc) CycleFocusForward() {
	w.cycleFocus(1)
}

func (w *Wtc) CycleFocusBackward() {
	w.cycleFocus(-1)
}

func (w *Wtc) cycleFocus(offset int) {
	var next int
	for i, panel := range w.panels {
		// NOTE: one (and only one) panel will always have a focus
		if panel.HasFocus() {
			next = abs(i + offset) % len(w.panels)
			break
		}
	}

	w.app.SetFocus(w.panels[next])
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
