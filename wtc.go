package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type TimeSnapshot struct {
	TotalSeconds int
}

func New(snapshot string) (*TimeSnapshot, error) {
	var hr, min, sec int

	if  m, err := regexp.MatchString(`^\d*$`, snapshot); m {
        if err != nil {
            return &TimeSnapshot{}, err
        }
		sec, err = strconv.Atoi(snapshot)
        if err != nil {
            return &TimeSnapshot{}, err
        }
	} else if  m, err := regexp.MatchString(`^\d+:\d{2}$`, snapshot); m {
        if err != nil {
            return &TimeSnapshot{}, err
        }
        s := strings.Split(snapshot, ":")
        min, err = strconv.Atoi(s[0])
        if err != nil {
            return &TimeSnapshot{}, err
        }
        sec, err = strconv.Atoi(s[1])
        if err != nil {
            return &TimeSnapshot{}, err
        }
        if sec >= 60 {
            return &TimeSnapshot{}, fmt.Errorf("seconds field must be less than 60")
        }
	} else if  m, err := regexp.MatchString(`^\d+:\d{2}:\d{2}$`, snapshot); m {
        if err != nil {
            return &TimeSnapshot{}, err
        }
        s := strings.Split(snapshot, ":")
        hr, err = strconv.Atoi(s[0])
        if err != nil {
            return &TimeSnapshot{}, err
        }
        min, err = strconv.Atoi(s[1])
        if err != nil {
            return &TimeSnapshot{}, err
        }
        sec, err = strconv.Atoi(s[2])
        if err != nil {
            return &TimeSnapshot{}, err
        }

        // Error handling
        if min >= 60 && sec >= 60 {
            err = fmt.Errorf("minutes and seconds")
        } else if min >= 60 {
            err = fmt.Errorf("minutes")
        } else if sec >= 60 {
            err = fmt.Errorf("seconds")
        }
        if err != nil {
            err = fmt.Errorf("%v field must be less than 60", err)
            return &TimeSnapshot{}, err
        }
    } else {
        return &TimeSnapshot{}, fmt.Errorf("incorrect format")
    }
    return &TimeSnapshot{ TotalSeconds: (hr * 3600) + (min * 60) + sec }, nil
}

func (t TimeSnapshot) String() string {
    hrs := t.TotalSeconds / 3600
    min := (t.TotalSeconds / 60) % 60
    sec := t.TotalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hrs, min, sec)
}

func countup(duration string) {
	for i := 0; ; i++ {
		fmt.Printf("%d\r", i)
		time.Sleep(time.Second)
	}
}

func countdown(duration string) {

}

func ArgParse() (string, string) {
    var duration string
    if len(os.Args) > 1 {
        duration = os.Args[1]
    }
	return duration, filepath.Base(os.Args[0])
}

func main() {
	// Handle interruptions
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		os.Exit(1)
	}()

	// Main program
	duration, prog := ArgParse()
    t, err := New(duration)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s: %v\n", prog, err)
        os.Exit(1)
    }
    fmt.Println(t)
	if duration == "" {
		// countup(duration)
	} else {
		// countdown(duration)
	}
}
