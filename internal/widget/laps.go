package widget

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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

// AddLap adds a new lap into l with Lap time lapSeconds and total time
// as totalSeconds.
func (l *LapTable) AddLap(lapSeconds int, totalSeconds int) *LapTable {
	// Accounting for header rows.
	l.InsertRow(2)

	var newCell = func(text string, ref interface{}) *tview.TableCell {
		c := tview.NewTableCell(text)
		c.SetReference(ref)
		c.SetAlign(tview.AlignCenter)
		c.SetStyle(tcell.StyleDefault.Background(l.GetBackgroundColor()).Foreground(tcell.ColorWhite))
		return c
	}

	// laps will start counting from 1.
	lap := l.GetRowCount() - 1
	l.SetCell(2, 0, newCell(fmt.Sprint(lap), lap))
	l.SetCell(2, 1, newCell(l.Format(lapSeconds), lapSeconds))
	l.SetCell(2, 2, newCell(l.Format(totalSeconds), totalSeconds))
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
