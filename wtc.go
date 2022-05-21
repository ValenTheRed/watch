package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type keyMapWtc struct {
	Quit, CycleFocusForward, CycleFocusBackward *Binding
}

func (km *keyMapWtc) Keys() []*Binding {
	return []*Binding{
		km.Quit, km.CycleFocusForward, km.CycleFocusBackward,
	}
}

type Paneler interface {
	tview.Primitive
	HasFocus() bool
}

type Wtc struct {
	app *tview.Application

	stopwatch *Stopwatch
	timer     *Timer
	help      *HelpView

	keyMap *keyMapWtc

	// panels is the list of widgets currently being displayed
	panels []Paneler
}

func NewWtc(app *tview.Application, duration int) *Wtc {
	w := &Wtc{
		app:  app,
		help: NewHelpView(),
		keyMap: &keyMapWtc{
			Quit: NewBinding(
				WithRune('q'), WithHelp("Quit"),
			),
			CycleFocusForward: NewBinding(
				WithKey(tcell.KeyTab), WithHelp("Cycle focus forward"),
			),
			CycleFocusBackward: NewBinding(
				WithKey(tcell.KeyBacktab), WithHelp("Cycle focus backward"),
			),
		},
	}

	w.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case w.keyMap.Quit.Rune():
				wtc.app.Stop()
			}
		case w.keyMap.CycleFocusForward.Key():
			wtc.CycleFocusForward()
		case w.keyMap.CycleFocusBackward.Key():
			wtc.CycleFocusBackward()
		}
		return event
	})

	w.InitMain(duration)
	w.help.SetGlobals(w.keyMap)

	return w
}

func (w *Wtc) InitMain(duration int) {
	var p Paneler
	if duration == 0 {
		w.stopwatch = NewStopwatch()
		p = w.stopwatch
	} else {
		w.timer = NewTimer(duration)
		p = w.timer
	}
	// w.help widget will not be a focus target.
	// See: [FIXME](utils.go: focusFunc())
	w.panels = append(w.panels, p)
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
	var next Paneler
	for i, panel := range w.panels {
		// NOTE: one (and only one) panel will always have a focus
		if panel.HasFocus() {
			next = w.panels[abs(i + offset) % len(w.panels)]
			break
		}
	}

	w.app.SetFocus(next)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
