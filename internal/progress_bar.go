package widget

import "github.com/rivo/tview"

type ProgressBar struct {
	*tview.Box

	// percent tracks the progress percentage, and belongs to the closed
	// interval [0, 100].
	percent int
}

// NewProgressBar returns a new ProgressBar initialised at 0% progress.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		Box: tview.NewBox(),
		percent: 0,
	}
}
