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

	return q
}
