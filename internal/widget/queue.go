package widget

import (
	"fmt"

	"github.com/rivo/tview"
)

type Queue struct {
	*Table

	// head is the currently selected row (after subtracting the
	// header rows).
	head int

	// An optional function which gets called whenever the user selects
	// a cell (eg: presses Enter on a cell). row is the row of the
	// selected cell. Row indexing starts with the row after the header
	// rows.
	selected func(row int)
}

const queueHeadIcon = "->"

// NewQueue returns a new Queue, with the duration column formatted
// using SecondWithColons.
func NewQueue(durations ...int) *Queue {
	q := &Queue{
		Table:  NewTable("Queue", "Timer duration"),
		head:   -1,
	}

	var newCell = func(text string, ref interface{}) *tview.TableCell {
		c := tview.NewTableCell(text)
		c.SetReference(ref)
		c.SetAlign(tview.AlignCenter)
		c.SetStyle(q.Table.GetCellStyle())
		return c
	}

	for i, duration := range durations {
		q.SetCell(i, 0, newCell(fmt.Sprint(i+1), i+1))
		q.SetCell(i, 1, newCell(SecondWithColons(duration), duration))
	}
	q.head = 0
	q.GetCell(0, 0).SetText(queueHeadIcon)

	// Pressing the Enter key leads to "selecting" that row.
	q.Table.SetSelectedFunc(func(row, column int) {
		// get row index after removing the header rows
		q.Select(row-2)
	})

	return q
}

// SetDurationFormat formats the duration column's text using format.
func (q *Queue) SetDurationFormat(format func(seconds int) string) *Queue {
	for r := 0; r < q.GetRowCount()-2; r++ {
		cell := q.GetCell(r, 1)
		cell.SetText(format(cell.GetReference().(int)))
	}
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

// Select selects row row. This also fires "selected" handler, if set.
// Row indexing starts with the row after the header rows.
func (q *Queue) Select(row int) *Queue {
	// Remove queueHeadIcon from current row.
	cell := q.GetCell(q.head, 0)
	cell.SetText(fmt.Sprint(cell.GetReference()))

	// Attach queueHeadIcon to the given row.
	q.head = row
	q.GetCell(q.head, 0).SetText(queueHeadIcon)

	if q.selected != nil {
		q.selected(row)
	}
	return q
}

// Next "selects" the next item from the queue.
func (q *Queue) Next() *Queue {
	next := q.head + 1
	// If queue has not reached it's last timer
	if next < q.GetRowCount() - 2 {
		q.Select(next)
	}
	return q
}

// Previous "selects" the previous item from the queue.
func (q *Queue) Previous() *Queue {
	prev := q.head - 1
	// If queue has not moved further than it's first timer
	if prev > -1 {
		q.Select(prev)
	}
	return q
}
