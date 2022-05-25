package stopwatch

import "github.com/rivo/tview"

// newLapCell returns a Table cell with a default style for a laps cell
// applied.
func newLapCell(text string, ref interface{}) *tview.TableCell {
	return tview.NewTableCell(text).
		SetReference(ref).
		SetAlign(tview.AlignCenter)
}
