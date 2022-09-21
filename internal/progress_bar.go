package widget

import "github.com/rivo/tview"

type ProgressBar struct {
	*tview.Box

	// percent tracks the progress percentage, and belongs to the closed
	// interval [0, 100].
	percent int
}
