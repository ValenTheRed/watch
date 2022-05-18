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

type TimerState int8

const (
	timerReset TimerState = 1 << iota
	timerPause
	timerStart
)

func main() {
	flag.Parse()
	duration, err := ParseDuration(flag.Arg(0))
	exitOnErr(err)

	stateMsg := make(chan TimerState)

	app := tview.NewApplication()
	textview := tview.NewTextView()
	textview.
		SetChangedFunc(func() {
			app.Draw()
		}).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'r':
				stateMsg <- timerReset
			case 'p':
				stateMsg <- timerPause
			case 's':
				stateMsg <- timerStart
			}
			return event
		}).
		SetBorder(true)

	if duration == 0 {
		go countup(textview, stateMsg)
	} else {
		go countdown(duration, stateMsg, textview)
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(textview, true).Run(); err != nil {
		panic(err)
	}
	close(stateMsg)
}

// countup counts up from zero seconds till infinity.
//
// Progress is written to tv.
func countup(tv *tview.TextView, msg <-chan TimerState) {
	tick := time.NewTicker(1 * time.Second)
	for t := 0; ; {
		select {
		case m := <-msg:
			switch m {
			case timerReset:
				t = 0
				tick.Reset(1 * time.Second)
			case timerPause:
				tick.Stop()
				for e := range msg {
					if e == timerStart {
						break
					}
				}
				tick.Reset(1 * time.Second)
			}
		case <-tick.C:
			tv.SetText(FormatSecond(t))
			t++
		}
	}
}

// countdown counts down for duration seconds, and runs pingFile when
// countdown ends.
//
// Progress is written to tv.
func countdown(duration int, msg <-chan TimerState, tv *tview.TextView) {
	tick := time.NewTicker(1 * time.Second)
	for t := duration; t > 0; {
		select {
		case m := <-msg:
			switch m {
			case timerReset:
				t = duration
				tick.Reset(1 * time.Second)
			case timerPause:
				tick.Stop()
				for e := range msg {
					if e == timerStart {
						break
					}
				}
				tick.Reset(1 * time.Second)
			}
		case <-tick.C:
			tv.SetText(FormatSecond(t))
			t--
		}
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
