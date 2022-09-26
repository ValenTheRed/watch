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
