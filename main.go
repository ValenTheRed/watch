// Program wtc implements a watch with timer and stopwatch functionality
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

const binaryName = "wtc"

var (
	usage = fmt.Sprintf("usage: %s [-help] [duration]", binaryName) + `
Terminal based watch with timer and stopwatch functionality.

Specify no arguments to start a stopwatch.
Specify duration to start a timer.

optional arguments:
duration	supported formats - [[hh:]mm:]ss
-help	    display this help message and exit`

	// Global controller for the whole application.
	wtc *Wtc
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}

	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", binaryName, err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	duration, err := ParseDuration(flag.Arg(0))
	exitOnErr(err)

	wtc = NewWtc(tview.NewApplication(), duration)

	if err := wtc.Run(); err != nil {
		panic(err)
	}
}
