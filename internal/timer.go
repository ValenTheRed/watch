package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Timer struct {
	*tview.Box

	// total is the total time in seconds.
	total int

	// elapsed is the time passed in seconds.
	elapsed int

	// Both determine the alignment of the timer text.
	verticalAlign, horizontalAlign int

	// stopCh will be used to signal to timer to stop ticking.
	stopCh chan struct{}

	// TextColor is the text color timer.
	TextColor tcell.Color

	// ShadowColor is the color for the shadow characters of the text.
	ShadowColor tcell.Color

	// done is an optional function that would be executed when timer
	// finishes.
	done func()
}

// NewTimer returns an initialised Timer that would counting down for
// duration seconds, and has it's text centered aligned both, vertically
// and horizontally.
func NewTimer(duration int) *Timer {
	return &Timer{
		Box:             tview.NewBox(),
		total:           duration,
		verticalAlign:   AlignCenter,
		horizontalAlign: tview.AlignCenter,
		// Timer will tick every second using Worker(). stopCh will be
		// used as the quit channel for Worker() and will be called from
		// the work() function being executed by Worker(). Since work()
		// is not executed in a go routine, signaling quit/stopCh
		// channel from work() would lead to a deadlock.
		stopCh:      make(chan struct{}, 1),
		TextColor:   tcell.ColorWhite,
		ShadowColor: tcell.ColorGrey,
	}
}

// SetDoneFunc sets a handler which is called when the timer has
// finished.
func (t *Timer) SetDoneFunc(handler func()) *Timer {
	t.done = handler
	return t
}

// IsTimeLeft returns whether t has count down for duration it was set
// for.
func (t *Timer) IsTimeLeft() bool {
	return t.elapsed < t.total
}

// Start starts the Timer if time is left.
func (t *Timer) Start() *Timer {
	if !t.IsTimeLeft() {
		return t
	}
	go Worker(func() {
		t.elapsed++
		if t.IsTimeLeft() {
			return
		}
		t.Stop()
		if t.done != nil {
			t.done()
		}
	}, t.stopCh)
	return t
}

// Stop stops the Timer if time is left.
func (t *Timer) Stop() *Timer {
	if t.IsTimeLeft() {
		t.stopCh <- struct{}{}
	}
	return t
}
