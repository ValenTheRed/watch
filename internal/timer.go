package widget

import "github.com/rivo/tview"

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
		stopCh: make(chan struct{}, 1),
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

// Stop stops the Timer if time is left.
func (t *Timer) Stop() *Timer {
	if t.IsTimeLeft() {
		t.stopCh <- struct{}{}
	}
	return t
}
