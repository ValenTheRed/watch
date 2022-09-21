package widget

import "time"

// Vertical alignment.
const (
	AlignCenter = iota
	AlignUp
	AlignDown
)

// ANSI Shadow font. Font character are guranteed to have a length of 6
// cells, and have all of the rows be equal in length.
// Ref: https://raw.githubusercontent.com/anonymoushack47/ANSI-shadow/master/ansi-shadow.flf
var ANSIShadow = map[rune][]string{
	'0': {
		" ██████╗ ",
		"██╔═████╗",
		"██║██╔██║",
		"████╔╝██║",
		"╚██████╔╝",
		" ╚═════╝ ",
	},
	'1': {
		" ██╗",
		"███║",
		"╚██║",
		" ██║",
		" ██║",
		" ╚═╝",
	},
	'2': {
		"██████╗ ",
		"╚════██╗",
		" █████╔╝",
		"██╔═══╝ ",
		"███████╗",
		"╚══════╝",
	},
	'3': {
		"██████╗ ",
		"╚════██╗",
		" █████╔╝",
		" ╚═══██╗",
		"██████╔╝",
		"╚═════╝ ",
	},
	'4': {
		"██╗  ██╗",
		"██║  ██║",
		"███████║",
		"╚════██║",
		"     ██║",
		"     ╚═╝",
	},
	'5': {
		"███████╗",
		"██╔════╝",
		"███████╗",
		"╚════██║",
		"███████║",
		"╚══════╝",
	},
	'6': {
		" ██████╗ ",
		"██╔════╝ ",
		"███████╗ ",
		"██╔═══██╗",
		"╚██████╔╝",
		" ╚═════╝ ",
	},
	'7': {
		"███████╗",
		"╚════██║",
		"    ██╔╝",
		"   ██╔╝ ",
		"   ██║  ",
		"   ╚═╝  ",
	},
	'8': {
		" █████╗ ",
		"██╔══██╗",
		"╚█████╔╝",
		"██╔══██╗",
		"╚█████╔╝",
		" ╚════╝ ",
	},
	'9': {
		" █████╗ ",
		"██╔══██╗",
		"╚██████║",
		" ╚═══██║",
		" █████╔╝",
		" ╚════╝ ",
	},
	':': {
		"   ",
		"██╗",
		"╚═╝",
		"██╗",
		"╚═╝",
		"   ",
	},
	'h': {
		"██╗  ██╗",
		"██║  ██║",
		"███████║",
		"██╔══██║",
		"██║  ██║",
		"╚═╝  ╚═╝",
	},
	'm': {
		"███╗   ███╗",
		"████╗ ████║",
		"██╔████╔██║",
		"██║╚██╔╝██║",
		"██║ ╚═╝ ██║",
		"╚═╝     ╚═╝",
	},
	's': {
		"███████╗",
		"██╔════╝",
		"███████╗",
		"╚════██║",
		"███████║",
		"╚══════╝",
	},
	' ': {
		"    ",
		"    ",
		"    ",
		"    ",
		"    ",
		"    ",
	},
}

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

func stringToANSIShadow(str string) []string {
	text := make([]string, 6)
	for i := 0; i < len(text); i++ {
		rs := []rune{}
		for _, r := range str {
			rs = append(rs, []rune(ANSIShadow[r][i])...)
		}
		text[i] = string(rs)
	}
	return text
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