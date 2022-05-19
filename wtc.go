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
	w := &Wtc{
		app:  app,
		main: tview.NewTextView(),
		help: tview.NewTextView(),
	}
	w.panels = []*tview.TextView{ w.main, w.help }

	for _, p := range w.panels {
		p := p
		p.
			SetTextAlign(tview.AlignCenter).
			SetTitleAlign(tview.AlignLeft).
			SetBorder(true).
			SetBackgroundColor(tcell.ColorDefault).
			SetFocusFunc(func() {
				p.
					SetTitleColor(tcell.ColorOrange).
					SetBorderColor(tcell.ColorOrange)
			})
			p.SetBlurFunc(func() {
				p.
					SetTitleColor(tview.Styles.TitleColor).
					SetBorderColor(tview.Styles.BorderColor)
			})
	}
	w.help.
		SetTitle("Help")

	w.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				wtc.app.Stop()
			}
		case tcell.KeyTab:
			wtc.CycleFocusForward()
		case tcell.KeyBacktab:
			wtc.CycleFocusBackward()
		}
		return event
	})

	return w
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
