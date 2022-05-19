package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Stopwatch struct {
	*tview.TextView
	elapsed int
	running bool
	stopMsg chan struct{}
}

func NewStopwatch() *Stopwatch {
	sw := &Stopwatch{
		TextView: tview.NewTextView(),
		stopMsg: make(chan struct{}),
	}
	sw.
		SetChangedFunc(func() {
			wtc.app.Draw()
		}).
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(func() {
			sw.
				SetTitleColor(tcell.ColorOrange).
				SetBorderColor(tcell.ColorOrange)
		}).
		SetBlurFunc(func() {
			sw.
				SetTitleColor(tview.Styles.TitleColor).
				SetBorderColor(tview.Styles.BorderColor)
		}).
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
		SetTitle("Stopwatch")

	sw.UpdateDisplay()
	return sw
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
