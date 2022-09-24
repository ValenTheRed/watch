package widget

type LapTable struct {
	*Table

	// Format will be used to format the lap and total time.
	Format func(seconds int) string
}
