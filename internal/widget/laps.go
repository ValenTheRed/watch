package widget

import (
	"fmt"

	"github.com/rivo/tview"
)

type LapTable struct {
	*Table

	// Format will be used to format the lap and total time.
	Format func(seconds int) string
}

// NewLapTable returns a new LapTable. The seconds are formatted using
// SecondWithColons by default.
func NewLapTable() *LapTable {
	t := NewTable("Lap", "Lap time", "Total")
	return &LapTable{
		Table:  t,
		Format: SecondWithColons,
	}
}

// AddLap adds a new lap into l with Lap time total time as
// totalSeconds.
func (l *LapTable) AddLap(totalSeconds int) *LapTable {
	var newCell = func(text string, ref interface{}) *tview.TableCell {
		c := tview.NewTableCell(text)
		c.SetReference(ref)
		c.SetAlign(tview.AlignCenter)
		c.SetStyle(l.Table.GetCellStyle())
		return c
	}

	var lap, lapSeconds int

	if l.GetRowCount() == 2 {
		lap, lapSeconds = 1, totalSeconds
	} else {
		i, _, total := l.GetLap(0)
		lap, lapSeconds = i+1, totalSeconds - total
	}

	l.InsertRow(2)
	l.SetCell(0, 0, newCell(fmt.Sprint(lap), lap))
	l.SetCell(0, 1, newCell(l.Format(lapSeconds), lapSeconds))
	l.SetCell(0, 2, newCell(l.Format(totalSeconds), totalSeconds))
	return l
}

// GetLap returns the lap at row row. Row indexing starts with the row
// after the header rows.
func (l *LapTable) GetLap(row int) (lap, lapSeconds, totalSeconds int) {
	var getData = func(col int) int {
		return l.GetCell(row, col).GetReference().(int)
	}
	return getData(0), getData(1), getData(2)
}

// GetHighlightedLap returns the currently highlighted lap.
func (l *LapTable) GetHighlightedLap() (lap, lapSeconds, totalSeconds int) {
	row, _ := l.GetSelection()
	return l.GetLap(row - 2)
}
