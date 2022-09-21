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

// Percent returns the progress percent.
func (p *ProgressBar) Percent() int {
	return p.percent
}

// SetPercent sets the progress to v percent. v must belong to the
// closed interval [0, 100].
func (p *ProgressBar) SetPercent(v int) *ProgressBar {
	p.percent = v
	return p
}
