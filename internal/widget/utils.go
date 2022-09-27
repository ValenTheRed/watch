package widget

import (
	"fmt"
	"strings"
	"time"
)

// Vertical alignment.
const (
	AlignCenter = iota
	AlignUp
	AlignDown
)

// ANSI Shadow font. Font character are guranteed to have a length of 6
// cells, and have all of the rows be equal in length.
// Ref: https://raw.githubusercontent.com/anonymoushack47/ANSI-shadow/master/ansi-shadow.flf
var ANSIShadow = map[rune][][]rune{
	'0': {
		[]rune(" ██████╗ "),
		[]rune("██╔═████╗"),
		[]rune("██║██╔██║"),
		[]rune("████╔╝██║"),
		[]rune("╚██████╔╝"),
		[]rune(" ╚═════╝ "),
	},
	'1': {
		[]rune(" ██╗"),
		[]rune("███║"),
		[]rune("╚██║"),
		[]rune(" ██║"),
		[]rune(" ██║"),
		[]rune(" ╚═╝"),
	},
	'2': {
		[]rune("██████╗ "),
		[]rune("╚════██╗"),
		[]rune(" █████╔╝"),
		[]rune("██╔═══╝ "),
		[]rune("███████╗"),
		[]rune("╚══════╝"),
	},
	'3': {
		[]rune("██████╗ "),
		[]rune("╚════██╗"),
		[]rune(" █████╔╝"),
		[]rune(" ╚═══██╗"),
		[]rune("██████╔╝"),
		[]rune("╚═════╝ "),
	},
	'4': {
		[]rune("██╗  ██╗"),
		[]rune("██║  ██║"),
		[]rune("███████║"),
		[]rune("╚════██║"),
		[]rune("     ██║"),
		[]rune("     ╚═╝"),
	},
	'5': {
		[]rune("███████╗"),
		[]rune("██╔════╝"),
		[]rune("███████╗"),
		[]rune("╚════██║"),
		[]rune("███████║"),
		[]rune("╚══════╝"),
	},
	'6': {
		[]rune(" ██████╗ "),
		[]rune("██╔════╝ "),
		[]rune("███████╗ "),
		[]rune("██╔═══██╗"),
		[]rune("╚██████╔╝"),
		[]rune(" ╚═════╝ "),
	},
	'7': {
		[]rune("███████╗"),
		[]rune("╚════██║"),
		[]rune("    ██╔╝"),
		[]rune("   ██╔╝ "),
		[]rune("   ██║  "),
		[]rune("   ╚═╝  "),
	},
	'8': {
		[]rune(" █████╗ "),
		[]rune("██╔══██╗"),
		[]rune("╚█████╔╝"),
		[]rune("██╔══██╗"),
		[]rune("╚█████╔╝"),
		[]rune(" ╚════╝ "),
	},
	'9': {
		[]rune(" █████╗ "),
		[]rune("██╔══██╗"),
		[]rune("╚██████║"),
		[]rune(" ╚═══██║"),
		[]rune(" █████╔╝"),
		[]rune(" ╚════╝ "),
	},
	':': {
		[]rune("   "),
		[]rune("██╗"),
		[]rune("╚═╝"),
		[]rune("██╗"),
		[]rune("╚═╝"),
		[]rune("   "),
	},
	'h': {
		[]rune("██╗  ██╗"),
		[]rune("██║  ██║"),
		[]rune("███████║"),
		[]rune("██╔══██║"),
		[]rune("██║  ██║"),
		[]rune("╚═╝  ╚═╝"),
	},
	'm': {
		[]rune("███╗   ███╗"),
		[]rune("████╗ ████║"),
		[]rune("██╔████╔██║"),
		[]rune("██║╚██╔╝██║"),
		[]rune("██║ ╚═╝ ██║"),
		[]rune("╚═╝     ╚═╝"),
	},
	's': {
		[]rune("███████╗"),
		[]rune("██╔════╝"),
		[]rune("███████╗"),
		[]rune("╚════██║"),
		[]rune("███████║"),
		[]rune("╚══════╝"),
	},
	' ': {
		[]rune("    "),
		[]rune("    "),
		[]rune("    "),
		[]rune("    "),
		[]rune("    "),
		[]rune("    "),
	},
	'-': {
		[]rune("       "),
		[]rune("       "),
		[]rune("██████╗"),
		[]rune("╚═════╝"),
		[]rune("       "),
		[]rune("       "),
	},
}

// getCenter returns the coordinate from where, if drawn, an object of
// length reservedLen would looked centered on a screen of size
// totalLen. Assuming 0 is the origin.
// Eg.
// total = 100, reserved = 10 => center = (100-10)/2 = floor(45) = 45
// So, in the interval [0, 99], if drawn from point 45, an object of
// length 10 will look centered.
func getCenter(totalLen, reservedLen int) int {
	return (totalLen - reservedLen) / 2
}

// decomposeSecond breaks seconds s into hours, minutes and seconds.
func DecomposeSecond(s int) (hrs, min, sec int) {
	return s / 3600, (s / 60) % 60, s % 60
}

// SecondWithLetters formats seconds s as 'XXh XXm XXs' or 'XXm XXs' or
// 'XXs'. Leading zeros are omitted.
func SecondWithLetters(s int) string {
	hrs, min, sec := DecomposeSecond(s)
	var str strings.Builder
	if hrs != 0 {
		str.WriteString(fmt.Sprintf("%dh ", hrs))
	}
	if hrs != 0 || min != 0 {
		str.WriteString(fmt.Sprintf("%dm ", min))
	}
	str.WriteString(fmt.Sprintf("%ds", sec))
	return str.String()
}

// SecondWithColons formats seconds s as 'XX:XX:XX' or 'XX:XX'. Leading
// zeros are not omitted.
func SecondWithColons(s int) string {
	hrs, min, sec := DecomposeSecond(s)
	var str strings.Builder
	if hrs != 0 {
		str.WriteString(fmt.Sprintf( "%02d:", hrs))
	}
	str.WriteString(fmt.Sprintf("%02d:%02d", min, sec))
	return str.String()
}

// SecondToANSIShadowWithLetters returns s in the format, 12h 34m 55s,
// in ANSIShadow font. If hours in zero, then minutes will be omitted if
// it is zero. If hours is not zero, minutes is not omitted.
func SecondToANSIShadowWithLetters(s int) []string {
	return stringToANSIShadow(SecondWithLetters(s))
}

// SecondToANSIShadowWithColons formats seconds s as 'XX:XX:XX' or
// 'XX:XX', without omitting leading zeros, and returns it in ANSI
// Shadow font.
func SecondToANSIShadowWithColons(s int) []string {
	return stringToANSIShadow(SecondWithColons(s))
}

func stringToANSIShadow(str string) []string {
	text := make([]string, 6)
	for i := 0; i < len(text); i++ {
		rs := []rune{}
		for _, r := range str {
			rs = append(rs, ANSIShadow[r][i]...)
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
