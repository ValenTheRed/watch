// Program wtc implements a watch with timer and stopwatch functionality
package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
    "flag"
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

type TimeSnapshot struct {
	TotalSeconds int
}

// Return error if sec/min field are not less than 60
func CheckField(sec, min int, checkMin bool) error {
	var err error
	if checkMin && min >= 60 && sec >= 60 {
		err = fmt.Errorf("minute's and second's")
	} else if checkMin && min >= 60 {
		err = fmt.Errorf("minute's")
	} else if sec >= 60 {
		err = fmt.Errorf("second's")
	}
	if err != nil {
		err = fmt.Errorf("%v field must be less than 60", err)
	}
	return err
}

// Returns a new TimeSnapshot whose TotalSeconds is equal to the
// total seconds in some snapshot of format [[hh:]mm:]ss
func New(snapshot string) (*TimeSnapshot, error) {
	var hr, min, sec int

	if m, err := regexp.MatchString(`^\d*$`, snapshot); m {
		if err != nil {
			return &TimeSnapshot{}, err
		}
		sec, _ = strconv.Atoi(snapshot)
	} else if m, err := regexp.MatchString(`^\d+:\d{2}$`, snapshot); m {
		if err != nil {
			return &TimeSnapshot{}, err
		}
		s := strings.Split(snapshot, ":")
		min, _ = strconv.Atoi(s[0])
		sec, _ = strconv.Atoi(s[1])
		if err = CheckField(sec, min, false); err != nil {
			return &TimeSnapshot{}, err
		}
	} else if m, err := regexp.MatchString(`^\d+:\d{2}:\d{2}$`, snapshot); m {
		if err != nil {
			return &TimeSnapshot{}, err
		}
		s := strings.Split(snapshot, ":")
		hr, _ = strconv.Atoi(s[0])
		min, _ = strconv.Atoi(s[1])
		sec, _ = strconv.Atoi(s[2])
		if err = CheckField(sec, min, true); err != nil {
			return &TimeSnapshot{}, err
		}
	} else {
        return &TimeSnapshot{}, fmt.Errorf("duration must be in [[hh:]mm:]ss format")
	}
	return &TimeSnapshot{TotalSeconds: (hr * 3600) + (min * 60) + sec}, nil
}

func (t TimeSnapshot) String() string {
	hrs := t.TotalSeconds / 3600
	min := (t.TotalSeconds / 60) % 60
	sec := t.TotalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hrs, min, sec)
}

// Countup from some time instant t till infinity
func (t TimeSnapshot) Countup() {
	for rsec := t.TotalSeconds; ; rsec++ {
		fmt.Printf("%s\r", TimeSnapshot{TotalSeconds: rsec})
		time.Sleep(1 * time.Second)
	}
}

// Countdown from some time instant t till zero seconds
func (t TimeSnapshot) Countdown() {
	for rsec := t.TotalSeconds; rsec > 0; rsec-- {
		fmt.Printf("%s\r", TimeSnapshot{TotalSeconds: rsec})
		time.Sleep(1 * time.Second)
	}
}

func ArgParse() string {
    flagHelp := flag.Bool("help", false, "display this help message and exit")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "%s\n", progUsage)
    }
    flag.Parse()

    if *flagHelp {
        fmt.Println(usage)
        os.Exit(0)
    }
    duration := flag.Arg(0)
	return duration
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

	duration := ArgParse()
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
