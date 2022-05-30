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

func NewWtc(app *tview.Application, durations []int) *Wtc {
	w := &Wtc{
		app:  app,
		help: help.NewHelpView(app),
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
	w.help.SetGlobals(w)
	w.help.UpdateDisplay()

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

	// w.help widget will not be a focus target, so will not be included
	// in w.panels.
	if len(durations) == 0 {
		w.stopwatch = stopwatch.New(w.app)
		w.stopwatch.Swtc.
			SetFocusFunc(focusFunc(w.stopwatch.Swtc, w.stopwatch)).
			SetBlurFunc(blurFunc(w.stopwatch.Swtc))
		w.stopwatch.Laps.
			SetSelectedStyle(
				tcell.Style{}.
					Background(tcell.NewHexColor(0xc591e9)).
					Foreground(tview.Styles.PrimitiveBackgroundColor),
			).
			SetFocusFunc(focusFunc(w.stopwatch.Laps, w.stopwatch.Laps)).
			SetBlurFunc(blurFunc(w.stopwatch.Laps))
		w.panels = []Paneler{w.stopwatch.Swtc, w.stopwatch.Laps}
	} else {
		w.timer = timer.New(durations, w.app)
		w.timer.Timer.
			SetFocusFunc(focusFunc(w.timer.Timer, w.timer)).
			SetBlurFunc(blurFunc(w.timer.Timer))
		w.timer.Queue.
			SetSelectedStyle(
				tcell.Style{}.
					Background(tcell.NewHexColor(0xc591e9)).
					Foreground(tview.Styles.PrimitiveBackgroundColor),
			).
			SetFocusFunc(focusFunc(w.timer.Queue, w.timer.Queue)).
			SetBlurFunc(blurFunc(w.timer.Queue))
		w.panels = []Paneler{w.timer.Timer, w.timer.Queue}
	}

	return w
}

func (w *Wtc) Keys() []*help.Binding {
	return []*help.Binding{
		w.km.Quit, //w.km.CycleFocusForward, w.km.CycleFocusBackward,
	}
}

func (w *Wtc) Run() error {
	if w.timer != nil {
		w.timer.Start()
	} else {
		w.stopwatch.Start()
	}
	return w.app.SetRoot(w.setLayout(), true).EnableMouse(true).Run()
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
		AddItem(w.help, 2, 0, false)
}

func (w *Wtc) CycleFocusForward() {
	w.cycleFocus(1)
}

func (w *Wtc) CycleFocusBackward() {
	w.cycleFocus(-1)
}

func (w *Wtc) cycleFocus(offset int) {
	if len(w.panels) == 1 && w.panels[0].HasFocus() {
		return
	}

	// Since `w.help` isn't included in the focusable group, with mouse
	// enabled, `w.help` can be brought to focus. In this case, `next`
	// will be nil.
	next := w.panels[0]
	for i, panel := range w.panels {
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
