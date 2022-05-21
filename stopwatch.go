package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type keyMapStopwatch struct {
	Reset, Stop, Start *Binding
}

func (km keyMapStopwatch) Keys() []*Binding {
	return []*Binding{km.Reset, km.Start, km.Stop}
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
			Reset: NewBinding(
				WithRune('r'), WithHelp("Reset"),
			),
			Stop: NewBinding(
				WithRune('p'), WithHelp("Pause"),
			),
			Start: NewBinding(
				WithRune('s'), WithHelp("Start"),
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
			case 'r':
				sw.Reset()
			case 'p':
				sw.Stop()
			case 's':
				sw.Start()
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
	sw.SetText(FormatSecond(sw.elapsed))
}

func (sw *Stopwatch) Start() {
	if !sw.running {
		sw.running = true
		go worker(func() {
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
