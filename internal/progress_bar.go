package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/rivo/tview"
)

type ProgressBar struct {
	*tview.Box

	// percent tracks the progress percentage, and belongs to the closed
	// interval [0, 100].
	percent int

	// align determines the vertical alignment of the progress bar.
	//
	// There is no horizontal alignment since the progress bar will fill
	// all of the available width.
	align int

	// color is the color for the fill characters of the progress bar.
	color tcell.Color

	// shadowColor is the color for the shadow characters of the
	// progress bar.
	shadowColor tcell.Color
}

// NewProgressBar returns a new ProgressBar initialised at 0% progress
// and center aligned.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		Box: tview.NewBox(),
		percent: 0,
		align: AlignCenter,
		color: tcell.ColorWhite,
		shadowColor: tcell.ColorGrey,
	}
}

// SetColor sets the color for the fill characters.
func (p *ProgressBar) SetColor(c colorful.Color) *ProgressBar {
	p.color = tcell.GetColor(c.Hex())
	return p
}

// SetColor sets the color for the shadow characters.
func (p *ProgressBar) SetShadowColor(c colorful.Color) *ProgressBar {
	p.shadowColor = tcell.GetColor(c.Hex())
	return p
}

// SetAlign sets the vertical alignment of the progress bar. Must be one
// of AlignCenter, AlignDown or AlignUp.
func (p *ProgressBar) SetAlign(align int) *ProgressBar {
	p.align = align
	return p
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
