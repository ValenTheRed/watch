// Program wtc implements a watch with timer and stopwatch functionality
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
)

const binaryName = "wtc"

var (
	//go:embed "ping.flac"
	pingFile []byte
	usage = fmt.Sprintf("usage: %s [-help] [duration]", binaryName) + `
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

func main() {
	// Handle signals SIGINT and SIGTERM
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		os.Exit(1)
	}()

	flag.Parse()

	if err := run(flag.Arg(0)); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", binaryName, err)
		os.Exit(1)
	}
}

func run(duration string) error {
	t, err := ParseDuration(duration)
	if err != nil {
		return err
	}

	if t == 0 {
		Countup(t)
	} else {
		Countdown(t)
		go fmt.Println("Time's up!")

		// Ping will close pingFile
		if err := Ping(bytes.NewReader(pingFile)); err != nil {
			return err
		}
	}

	return nil
}

func Ping(r io.Reader) error {
	streamer, format, err := flac.Decode(r)
	if err != nil {
		return err
	}
	defer streamer.Close()

	done := make(chan struct{})
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second / 10))
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- struct{}{}
	})))

	<-done
	return nil
}
