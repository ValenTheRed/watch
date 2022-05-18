// Program wtc implements a watch with timer and stopwatch functionality
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
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

	wtc = NewWtc(tview.NewApplication())

	if duration == 0 {
		sw := NewStopwatch()
		sw.Start()
	} else {
		t := NewTimer(duration)
		t.Start()
	}

	wtc.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			wtc.app.Stop()
		}
		return event
	})

	if err := wtc.app.SetRoot(wtc.main, true).Run(); err != nil {
		panic(err)
	}
}
