package widget

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
		Table: t,
		Format: SecondWithColons,
	}
}
