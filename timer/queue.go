package timer

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

type queue struct {
	*tview.Table
	km    map[string]*help.Binding
	title string
	// head tracks durations, not rows so, will always be one less than
	// rows.
	head int
}

// newQueue returns a new queue.
func newQueue() *queue {
	return &queue{
		Table: tview.NewTable(),
		title: " Queue ",
		km: map[string]*help.Binding{
			"Select": help.NewBinding(
				help.WithKey(tcell.KeyEnter),
				help.WithHelp("Select"),
			),
		},
	}
}

// init returns an initialised q. Should be run immediately after
// newQueue().
func (q *queue) init() *queue {
	q.
		initFirstRow().
		// column headers will always remain in view
		SetFixed(1, 0).
		SetSelectable(true, true).
		// select the first duration cell
		Select(1, 1).
		SetWrapSelection(true, false).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(q.title)

	// Switch off default keybinds for moving between columns.
	q.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h', 'l':
				return nil
			}
		case tcell.KeyLeft, tcell.KeyRight:
			return nil
		}
		return event
	})

	return q
}

// initFirstRow inserts and initialises the first row i.e. the row with
// column headers, for q.
func (q *queue) initFirstRow() *queue {
	if q.GetRowCount() > 0 {
		return q
	}
	q.InsertRow(0)
	// SetExpansion applies on the whole column.
	q.SetCell(0, 0,
		newQueueCell("Turn", nil).
			SetAttributes(tcell.AttrBold).
			SetExpansion(1),
	)
	q.SetCell(0, 1,
		newQueueCell("Duration", nil).
			SetAttributes(tcell.AttrBold).
			SetExpansion(2),
	)
	return q
}

// Keys returns the key bound to q.
func (q *queue) Keys() []*help.Binding {
	return []*help.Binding{
		q.km["Select"],
	}
}

// Title returns the title of q.
func (q *queue) Title() string {
	return q.title
}

// addDuration adds a new entry to q.
func (q *queue) addDuration(d int) *queue {
	// Table automatically adds the required cells without having to
	// insert a row first.
	row := q.GetRowCount()
	q.SetCell(row, 0, newQueueCell(fmt.Sprint(row), row))
	q.SetCell(row, 1, newQueueCell(utils.FormatSecond(d), d))
	return q
}

// getCurrentDuration returns the duration of the current head of the
// queue.
func (q *queue) getCurrentDuration() int {
	cell := q.GetCell(q.head+1, 1)
	return cell.Reference.(int)
}

// newQueueCell returns a Table cell with a default style for a laps cell
// applied.
func newQueueCell(text string, ref interface{}) *tview.TableCell {
	return tview.NewTableCell(text).
		SetReference(ref).
		SetAlign(tview.AlignCenter)
}
