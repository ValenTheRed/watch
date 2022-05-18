package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type stopwatch struct {
	*tview.TextView
	elapsed int
	running bool
	stopMsg chan struct{}
}

func NewStopwatch() *stopwatch {
	sw := &stopwatch{
		TextView: tview.NewTextView(),
		stopMsg:  make(chan struct{}),
	}
	sw.TextView.
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
		})
	return sw
}

func (sw *stopwatch) UpdateDisplay() {
	sw.SetText(FormatSecond(sw.elapsed))
}

func (sw *stopwatch) Start() {
	if !sw.running {
		sw.running = true
		go worker(func() {
			sw.UpdateDisplay()
			sw.elapsed++
		}, sw.stopMsg)
	}
}

func (sw *stopwatch) Stop() {
	if sw.running {
		sw.running = false
		sw.stopMsg <- struct{}{}
	}
}

func (sw *stopwatch) Reset() {
	sw.Stop()
	sw.elapsed = 0
	sw.UpdateDisplay()
	sw.Start()
}
