package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/stopwatch"
	"github.com/ValenTheRed/watch/timer"
)

type keyMap struct {
	Quit, CycleFocusForward, CycleFocusBackward *help.Binding
}

type Paneler interface {
	tview.Primitive
	HasFocus() bool
}

type Wtc struct {
	app *tview.Application

	stopwatch *stopwatch.Stopwatch
	timer     *timer.Timer
	help      *help.HelpView

	km keyMap

	// panels is the list of widgets currently being displayed
	panels []Paneler
}

func NewWtc(app *tview.Application, duration int) *Wtc {
	w := &Wtc{
		app:  app,
		help: help.NewHelpView(),
		km: keyMap{
			Quit: help.NewBinding(
				help.WithRune('q'), help.WithHelp("Quit"),
			),
			CycleFocusForward: help.NewBinding(
				help.WithKey(tcell.KeyTab),
				help.WithHelp("Cycle focus forward"),
			),
			CycleFocusBackward: help.NewBinding(
				help.WithKey(tcell.KeyBacktab),
				help.WithHelp("Cycle focus backward"),
			),
		},
	}

	w.help.
		SetFocusFunc(focusFunc(w.help, w.help)).
		SetBlurFunc(blurFunc(w.help))

	w.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case w.km.Quit.Rune():
				wtc.app.Stop()
			}
		case w.km.CycleFocusForward.Key():
			wtc.CycleFocusForward()
		case w.km.CycleFocusBackward.Key():
			wtc.CycleFocusBackward()
		}
		return event
	})

	w.InitMain(duration)
	w.help.SetGlobals(w)

	return w
}

func (w *Wtc) Keys() []*help.Binding {
	return []*help.Binding{
		w.km.Quit, w.km.CycleFocusForward, w.km.CycleFocusBackward,
	}
}

func (w *Wtc) InitMain(duration int) {
	var p Paneler
	if duration == 0 {
		w.stopwatch = stopwatch.New()
		p = w.stopwatch
		w.stopwatch.
			SetFocusFunc(focusFunc(w.stopwatch, w.stopwatch)).
			SetBlurFunc(blurFunc(w.stopwatch))
	} else {
		w.timer = timer.New(duration)
		p = w.timer
		w.timer.
			SetFocusFunc(focusFunc(w.timer, w.timer)).
			SetBlurFunc(blurFunc(w.timer))
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
			next = w.panels[abs(i+offset)%len(w.panels)]
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
