package widget

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	*tview.Table

	// The default style for data cells.
	cellStyle tcell.Style
}

// NewTable returns a new Table.
func NewTable(headers ...string) *Table {
	var newCell = func(text string, style tcell.Style) *tview.TableCell {
		c := tview.NewTableCell(text)
		c.SetStyle(style)
		c.SetAlign(tview.AlignCenter)
		c.SetExpansion(1)
		c.NotSelectable = true
		return c
	}

	t := tview.NewTable()
	defStyle := tcell.StyleDefault
	for i, header := range headers {
		t.SetCell(0, i, newCell(header, defStyle))
	}
	for i, header := range headers {
		t.SetCell(1, i, newCell( strings.Repeat("▔", len(header)), defStyle))
	}

	t.SetSelectable(true, false)
	t.SetFixed(2, 0)
	return &Table{t, defStyle}
}

// GetCell returns the cell at the given position. The position
// calculation doesn't consider the header cells. A row of 0 and col of
// 0 returns the first cell with data.
//
// A valid TableCell object is always returned but it will be
// uninitialized if the cell was not previously set. Such an
// uninitialized object will not automatically be inserted. Therefore,
// repeated calls to this function may return different pointers for
// uninitialized cells.
func (t *Table) GetCell(row, col int) *tview.TableCell {
	return t.Table.GetCell(row+2, col)
}

// SetCell sets a cell at the given position. The position calculation
// doesn't consider the header cells. A row of 0 and col of 0 returns
// the first cell with data.
//
// It is ok to directly instantiate a TableCell object. If the cell has
// content, at least the Text and Color fields should be set.
//
// Note that setting cells in previously unknown rows and columns will
// automatically extend the internal table representation with empty
// TableCell objects, e.g. starting with a row of 100,000 will
// immediately create 100,000 empty rows.
//
// To avoid unnecessary garbage collection, fill columns from left to
// right.
func (t *Table) SetCell(row, col int, cell *tview.TableCell) *Table {
	t.Table.SetCell(row+2, col, cell)
	return t
}

// SetHeaderStyle sets the style of all header cells as s.
func (t *Table) SetHeaderStyle(s tcell.Style) *Table {
	for i := 0; i < t.GetColumnCount(); i++ {
		t.Table.GetCell(0, i).SetStyle(s)
	}
	return t
}

// SetUnderlineStyle sets the style for all of the underline cells as s.
func (t *Table) SetUnderlineStyle(s tcell.Style) *Table {
	for i := 0; i < t.GetColumnCount(); i++ {
		t.Table.GetCell(1, i).SetStyle(s)
	}
	return t
}

// GetCellStyle returns Table's cell style.
func (t *Table) GetCellStyle() tcell.Style {
	return t.cellStyle
}

func (t *Table) SetCellStyle(s tcell.Style) *Table {
	for r := 0; r < t.GetRowCount(); r++ {
		for c := 0; c < t.GetColumnCount(); c++ {
			t.GetCell(r, c).SetStyle(s)
		}
	}
	return t
}
