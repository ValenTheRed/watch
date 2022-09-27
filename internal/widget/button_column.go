package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

// ButtonColumn draws contained tview.Button in a columnar fashion,
// like so (the borders may or may not be present):
// | Button 1 | | Button 2 | ... | Button N |
type ButtonColumn struct {
	*tview.Box

	// The contained buttons.
	buttons []*tview.Button

	// Both determine the alignment of the buttons.
	horizontalAlign, verticalAlign int
}

// NewButtonColumn returns a new ButtonColumn. It also set the left and
// right border padding of the buttons to 1.
func NewButtonColumn(buttons []*tview.Button) *ButtonColumn {
	for _, b := range buttons {
		b.SetBorderPadding(0, 0, 1, 1)
	}
	return &ButtonColumn{
		Box:             tview.NewBox(),
		buttons:         buttons,
		horizontalAlign: tview.AlignCenter,
		verticalAlign:   AlignCenter,
	}
}

// SetHorizontalAlign sets the veritcal alignment of the buttons. Must be
// one of tview.AlignCenter, tview.AlignLeft or tview.AlignRight.
func (br *ButtonColumn) SetHorizontalAlign(align int) *ButtonColumn {
	br.horizontalAlign = align
	return br
}

// SetVerticalAlign sets the veritcal alignment of the buttons. Must be
// one of AlignCenter, AlignUp or AlignDown.
func (br *ButtonColumn) SetVerticalAlign(align int) *ButtonColumn {
	br.verticalAlign = align
	return br
}

// HasFocus returns whether or not this primitive has focus.
func (bc *ButtonColumn) HasFocus() bool {
	for _, b := range bc.buttons {
		if b.HasFocus() {
			return true
		}
	}
	return bc.Box.HasFocus()
}

func (bc *ButtonColumn) Draw(screen tcell.Screen) {
	bc.DrawForSubclass(screen, bc)

	maxButtonWidth := func() int {
		m := -1
		for _, b := range bc.buttons {
			// 2 account for the left and right padding of the buttons.
			m = max(m, 2+runewidth.StringWidth(b.GetLabel()))
		}
		return m
	}()

	x, y, width, height := bc.GetInnerRect()

	if fixedHeight := 1; bc.verticalAlign == AlignCenter {
		y += getCenter(height, fixedHeight)
	} else if bc.verticalAlign == AlignDown {
		y += height - fixedHeight
	}
	if fixedWidth := maxButtonWidth*len(bc.buttons) + (len(bc.buttons) - 1); bc.horizontalAlign == tview.AlignCenter {
		x += getCenter(width, fixedWidth)
	} else if bc.horizontalAlign == tview.AlignRight {
		x += width - fixedWidth
	}

	sepStyle := tcell.StyleDefault.Background(bc.GetBackgroundColor())
	for _, b := range bc.buttons {
		b.SetRect(x, y, maxButtonWidth, 1)
		b.Draw(screen)
		x += maxButtonWidth
		screen.SetContent(x, y, ' ', nil, sepStyle)
		x += 1
	}
}

// MouseHandler returns the mouse handler for this primitive.
func (bc *ButtonColumn) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(p tview.Primitive)) (bool, tview.Primitive) {
	return bc.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !bc.InRect(event.Position()) {
			return false, nil
		}
		// Pass mouse events along to the first button that takes it.
		for _, b := range bc.buttons {
			if b == nil {
				continue
			}
			consumed, capture = b.MouseHandler()(action, event, setFocus)
			if consumed {
				return
			}
		}
		return
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
