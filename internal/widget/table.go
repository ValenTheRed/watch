package widget

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	*tview.Table
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
	defStyle := tcell.StyleDefault.Background(t.GetBackgroundColor())
	for i, header := range headers {
		t.SetCell(0, i, newCell(
			header,
			defStyle.Foreground(tcell.ColorWhite),
		))
	}
	for i, header := range headers {
		t.SetCell(1, i, newCell(
			strings.Repeat("▔", len(header)),
			defStyle.Foreground(tcell.ColorBlue),
		))
	}

	t.SetSelectable(true, false)
	return &Table{t}
}
