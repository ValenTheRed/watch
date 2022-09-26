package widget

import "github.com/rivo/tview"

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
