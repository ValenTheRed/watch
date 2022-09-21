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
