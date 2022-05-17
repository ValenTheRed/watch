// Program wtc implements a watch with timer and stopwatch functionality
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const binaryName = "wtc"

var (
	//go:embed "ping.flac"
	pingFile []byte
	usage    = fmt.Sprintf("usage: %s [-help] [duration]", binaryName) + `
Terminal based watch with timer and stopwatch functionality.

Specify no arguments to start a stopwatch.
Specify duration to start a timer.

optional arguments:
duration	supported formats - [[hh:]mm:]ss
-help	    display this help message and exit`
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

	app := tview.NewApplication()
	textview := tview.NewTextView()
	textview.
		SetChangedFunc(func() {
			app.Draw()
		}).
		SetBorder(true)

	if duration == 0 {
		go countup(textview)
	} else {
		go countdown(duration, textview)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		return nil
	})

	if err := app.SetRoot(textview, true).Run(); err != nil {
		panic(err)
	}
}

// countup counts up from zero seconds till infinity.
//
// Progress is written to tv.
func countup(tv *tview.TextView) {
	for t := 0; ; t++ {
		tv.SetText(FormatSecond(t))
		time.Sleep(1 * time.Second)
	}
}

// countdown counts down for duration seconds, and runs pingFile when
// countdown ends.
//
// Progress is written to tv.
func countdown(duration int, tv *tview.TextView) {
	for t := duration; t > 0; t-- {
		tv.SetText(FormatSecond(t))
		time.Sleep(1 * time.Second)
	}
	go tv.SetText(fmt.Sprintf("Your %s's up!\n", FormatSecond(duration)))

	// if err := Ping(bytes.NewReader(pingFile)); err != nil {
	//	 return err
	// }
	Ping(bytes.NewReader(pingFile))
}

func Ping(r io.Reader) error {
	streamer, format, err := flac.Decode(r)
	if err != nil {
		return err
	}
	defer streamer.Close()

	done := make(chan struct{})
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- struct{}{}
	})))

	<-done
	return nil
}
