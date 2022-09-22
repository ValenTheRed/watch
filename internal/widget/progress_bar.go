package widget

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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

	// TextColor is the color for the fill characters of the progress bar.
	TextColor tcell.Color

	// ShadowColor is the color for the shadow characters of the
	// progress bar.
	ShadowColor tcell.Color
}

// NewProgressBar returns a new ProgressBar initialised at 0% progress
// and center aligned.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		Box:         tview.NewBox(),
		percent:     0,
		align:       AlignCenter,
		TextColor:   tcell.ColorWhite,
		ShadowColor: tcell.ColorGrey,
	}
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
// closed interval [0, 100], returns p and error otherwise.
func (p *ProgressBar) SetPercent(v int) (*ProgressBar, error) {
	if v < 0 {
		return p, fmt.Errorf("progress: negative progress percent %d", v)
	} else if v > 100 {
		return p, fmt.Errorf("progress: progress percent %d larger than 100", v)
	}
	p.percent = v
	return p, nil
}

func (p *ProgressBar) Draw(screen tcell.Screen) {
	p.Box.DrawForSubclass(screen, p)

	// minHeight of the progress bar
	// ██
	// ██
	// ╚═
	const minHeight = 3

	x, y, width, height := p.GetInnerRect()
	if p.align == AlignCenter {
		y += getCenter(height, minHeight)
	} else if p.align == AlignDown {
		y += height - minHeight
	}

	const (
		fillChar             = '█'
		shadowHoriChar       = '═'
		shadowVertiChar      = '║'
		shadowUpperLeftChar  = '╔'
		shadowLowerLeftChar  = '╚'
		shadowUpperRightChar = '╗'
		shadowLowerRightChar = '╝'
	)

	shadowStyle := tcell.StyleDefault.Foreground(p.ShadowColor).Background(p.GetBackgroundColor())
	fillStyle := tcell.StyleDefault.Foreground(p.TextColor).Background(p.GetBackgroundColor())

	xProgressEnd := width * p.percent / 100
	xEnd := x + width - 1

	// line 1
	for i := 0; i < xProgressEnd; i++ {
		screen.SetCell(x+i, y, fillStyle, fillChar)
	}
	for i := xProgressEnd; i < width; i++ {
		screen.SetCell(x+i, y, shadowStyle, shadowHoriChar)
	}
	screen.SetCell(xEnd, y, shadowStyle, shadowUpperRightChar)
	if xProgressEnd == 0 {
		screen.SetCell(x, y, shadowStyle, shadowUpperLeftChar)
	}

	// line 2
	y++
	for i := 0; i < xProgressEnd; i++ {
		screen.SetCell(x+i, y, fillStyle, fillChar)
	}
	screen.SetCell(xEnd, y, shadowStyle, shadowVertiChar)
	if xProgressEnd == 0 {
		screen.SetCell(x, y, shadowStyle, shadowVertiChar)
	}

	// line 3
	y++
	screen.SetCell(x, y, shadowStyle, shadowLowerLeftChar)
	for i := 1; i < width; i++ {
		screen.SetCell(x+i, y, shadowStyle, shadowHoriChar)
	}
	screen.SetCell(xEnd, y, shadowStyle, shadowLowerRightChar)
}
