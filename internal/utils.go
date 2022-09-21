package widget

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
