package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ButtonColumn draws contained tview.Button in a columnar fashion,
// like so (the borders may or may not be present):
// | Button 1 | | Button 2 | ... | Button N |
type ButtonColumn struct {
	*tview.Box

	// The contained buttons.
	buttons []*tview.Button
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
