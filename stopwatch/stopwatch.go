package stopwatch

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

type keyMapStopwatch struct {
	Reset, Stop, Start *help.Binding
}

func (km keyMapStopwatch) Keys() []*help.Binding {
	return []*help.Binding{km.Reset, km.Start, km.Stop}
}

type Stopwatch struct {
	*tview.TextView
	elapsed int
	running bool
	stopMsg chan struct{}
	title   string
	keyMap  keyMapStopwatch
}

func NewStopwatch() *Stopwatch {
	sw := &Stopwatch{
		TextView: tview.NewTextView(),
		stopMsg:  make(chan struct{}),
		title:    " Stopwatch ",
		keyMap: keyMapStopwatch{
			Reset: help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			Stop: help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			Start: help.NewBinding(
				help.WithRune('s'),
				help.WithHelp("Start"),
				help.WithDisable(true),
			),
		},
	}

	sw.
		SetChangedFunc(func() {
			wtc.app.Draw()
		}).
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(focusFunc(sw, sw.keyMap)).
		SetBlurFunc(blurFunc(sw)).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case sw.keyMap.Reset.Rune():
				sw.Reset()
				sw.keyMap.Start.SetDisable(true)
				sw.keyMap.Stop.SetDisable(false)
			case sw.keyMap.Stop.Rune():
				if sw.keyMap.Stop.IsEnabled() {
					sw.Stop()
					sw.keyMap.Stop.SetDisable(true)
					sw.keyMap.Start.SetDisable(false)
				}
			case sw.keyMap.Start.Rune():
				if sw.keyMap.Start.IsEnabled() {
					sw.Start()
					sw.keyMap.Start.SetDisable(true)
					sw.keyMap.Stop.SetDisable(false)
				}
			}
			return event
		}).
		SetTitle(sw.title)

	sw.UpdateDisplay()
	return sw
}

func (sw *Stopwatch) Title() string {
	return sw.title
}

func (sw *Stopwatch) UpdateDisplay() {
	sw.SetText(utils.FormatSecond(sw.elapsed))
}

func (sw *Stopwatch) Start() {
	if !sw.running {
		sw.running = true
		go utils.Worker(func() {
			sw.elapsed++
			sw.UpdateDisplay()
		}, sw.stopMsg)
	}
}

func (sw *Stopwatch) Stop() {
	if sw.running {
		sw.running = false
		sw.stopMsg <- struct{}{}
	}
}

func (sw *Stopwatch) Reset() {
	sw.Stop()
	sw.elapsed = 0
	sw.UpdateDisplay()
	sw.Start()
}
