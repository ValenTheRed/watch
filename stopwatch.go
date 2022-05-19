package main

import (
	"github.com/gdamore/tcell/v2"
)

type Stopwatch struct {
	elapsed int
	running bool
	stopMsg chan struct{}
}

func NewStopwatch() *Stopwatch {
	sw := &Stopwatch{
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

func (sw *Stopwatch) UpdateDisplay() {
	wtc.main.SetText(FormatSecond(sw.elapsed))
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
