package widget

import "time"

// Vertical alignment.
const (
	AlignCenter = iota
	AlignUp
	AlignDown
)

// getCenter returns the coordinate from where, if drawn, an object of
// length reservedLen would looked centered on a screen of size
// totalLen. Assuming 0 is the origin.
// Eg.
// total = 100, reserved = 10 => center = (100-10-1)/2 = floor(44.5) = 44
// So, in the interval [0, 99], if drawn from point 44, an object of
// length 10 will look centered.
func getCenter(totalLen, reservedLen int) int {
	return (totalLen - reservedLen - 1) / 2
}

// decomposeSecond breaks seconds s into hours, minutes and seconds.
func DecomposeSecond(s int) (hrs, min, sec int) {
	return s / 3600, (s / 60) % 60, s % 60
}

// Worker executes work after every second. If a message is sent to
// quit, Worker returns.
func Worker(work func(), quit <-chan struct{}) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			work()
		case <-quit:
			return
		}
	}
}
