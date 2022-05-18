package main

import (
	"github.com/rivo/tview"
)

type timer struct {
	*tview.TextView
	duration, timeLeft int
	running            bool
	stopMsg            chan struct{}
}

func NewTimer(duration int) *timer {
	t := &timer{
		TextView: tview.NewTextView(),
		stopMsg:  make(chan struct{}),
		duration: duration,
		timeLeft: duration,
	}
	return t
}

func (t *timer) UpdateDisplay() {
	t.SetText(FormatSecond(t.timeLeft))
}

func (t *timer) Start() {
	if !t.running {
		t.running = true
		go worker(func() {
			if t.timeLeft > 0 {
				t.timeLeft--
				t.UpdateDisplay()
			} else {
				t.Stop()
			}
		}, t.stopMsg)
	}
}

func (t *timer) Stop() {
	if t.running {
		t.running = false
		t.stopMsg <- struct{}{}
	}
}

func (t *timer) Reset() {
	t.Stop()
	t.timeLeft = t.duration
	t.UpdateDisplay()
	t.Start()
}
