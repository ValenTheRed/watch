package widget

type Queue struct {
	*Table

	// head is the currently selected row (after subtracting the
	// header rows).
	head int

	// Format will be used to format seconds in durations column.
	Format func(seconds int) string
}
