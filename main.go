// Program wtc implements a watch with timer and stopwatch functionality
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
    progName  = "wtc"
    progUsage = "usage: wtc [-help] [duration]"
    usage     = progUsage + `
Terminal based watch with timer and stopwatch functionality.

Specify no arguments to start a stopwatch.
Specify duration to start a timer.

optional arguments:
  duration    supported formats - [[hh:]mm:]ss
  -help       display this help message and exit
    `
)

func init() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "%s\n", progUsage)
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

    duration := flag.Arg(0)
	t, err := New(duration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", progName, err)
		os.Exit(1)
	}
	if duration == "" {
		t.Countup()
	} else {
		t.Countdown()
		fmt.Println("Time's up!")
	}
}
