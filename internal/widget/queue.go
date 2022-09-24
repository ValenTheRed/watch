package widget

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Queue struct {
	*Table

	// head is the currently selected row (after subtracting the
	// header rows).
	head int

	// Format will be used to format seconds in durations column.
	Format func(seconds int) string

	// An optional function which gets called whenever the user selects
	// a cell (eg: presses Enter on a cell). row is the row of the
	// selected cell. Row indexing starts with the row after the header
	// rows.
	selected func(row int)
}

const queueHeadIcon = "->"

// NewQueue returns a new Queue.
func NewQueue(durations ...int) *Queue {
	q := &Queue{
		Table:  NewTable("Queue", "Timer duration"),
		head:   -1,
		Format: SecondWithColons,
	}

	var newCell = func(text string, ref interface{}) *tview.TableCell {
		c := tview.NewTableCell(text)
		c.SetReference(ref)
		c.SetAlign(tview.AlignCenter)
		c.SetStyle(tcell.StyleDefault.Background(q.GetBackgroundColor()).Foreground(tcell.ColorWhite))
		return c
	}

	for i, duration := range durations {
		q.SetCell(i, 0, newCell(fmt.Sprint(i+1), i+1))
		q.SetCell(i, 1, newCell(q.Format(duration), duration))
	}
	q.head = 0
	q.GetCell(0, 0).SetText(queueHeadIcon)

	q.Table.SetSelectedFunc(func(row, column int) {
		// remove the header rows from further calculations
		row -= 2

		// remove queueHeadIcon from the previously selected row.
		cell := q.GetCell(q.head, 0)
		cell.SetText(fmt.Sprint(cell.GetReference()))

		// attach queueHeadIcon to the currently selected row.
		q.head = row
		q.GetCell(q.head, 0).SetText(queueHeadIcon)

		if q.selected != nil {
			q.selected(row)
		}
	})

	return q
}

// SetSelectedFunc sets an optional function which gets called whenever
// the user selects a cell (eg: presses Enter on a cell). row is the row
// of the selected cell. Row indexing starts with the row after the
// header rows.
func (q *Queue) SetSelectedFunc(handler func(row int)) *Queue {
	q.selected = handler
	return q
}
