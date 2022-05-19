package main

import (
	"github.com/gdamore/tcell/v2"
)

type stopwatch struct {
	elapsed int
	running bool
	stopMsg chan struct{}
}

func NewStopwatch() *stopwatch {
	sw := &stopwatch{
		stopMsg: make(chan struct{}),
	}
	wtc.main.
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

func (sw *stopwatch) UpdateDisplay() {
	wtc.main.SetText(FormatSecond(sw.elapsed))
}

func (sw *stopwatch) Start() {
	if !sw.running {
		sw.running = true
		go worker(func() {
			sw.elapsed++
			sw.UpdateDisplay()
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
