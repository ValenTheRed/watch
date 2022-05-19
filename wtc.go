package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Panel interface {
	HasFocus() bool
}

type Wtc struct {
	app  *tview.Application

	stopwatch *Stopwatch
	timer *Timer
	help *tview.TextView

	// panels is the list of widgets currently being displayed
	panels []Panel
}

func NewWtc(app *tview.Application, duration int) *Wtc {
	w := &Wtc{
		app:  app,
		help: tview.NewTextView(),
	}

	w.help.
		SetBorder(true).
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

	w.InitMain(duration)

	return w
}

func (w *Wtc) InitMain(duration int) {
	var p Panel
	if duration == 0 {
		w.stopwatch = NewStopwatch()
		p = w.stopwatch
	} else {
		w.timer = NewTimer(duration)
		p = w.timer
	}
	w.panels = append(w.panels, p, w.help)
}

func (w *Wtc) Run() error {
	if w.timer != nil {
		w.timer.Start()
	} else {
		w.stopwatch.Start()
	}
	return w.app.SetRoot(w.setLayout(), true).Run()
}

func (w *Wtc) setLayout() *tview.Flex {
	var prim tview.Primitive
	if w.timer != nil {
		prim = w.timer
	} else {
		prim = w.stopwatch
	}
	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(prim, 0, 9, true).
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

	w.app.SetFocus(w.panels[next].(tview.Primitive))
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
