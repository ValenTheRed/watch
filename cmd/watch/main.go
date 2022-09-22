package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ValenTheRed/watch/internal/widget"
	"github.com/rivo/tview"
)

var (
	usage = `usage: watch [-help] [duration]
A clock with a stopwatch and a timer.

Specify a duration to start a timer. Or, leave it alone to start a stopwatch.

optional arguments:
duration    supported formats - [[hh:]mm:]ss
-help	    display this help message and exit`
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}
}

func main() {
	flag.Parse()

	durations := make([]int, len(flag.Args()))
	for i := range durations {
		var err error
		durations[i], err = ParseDuration(flag.Arg(i))
		if err != nil {
			log.Fatalln(fmt.Errorf("main: %v", err))
		}
		if durations[i] == 0 {
			log.Fatalln(fmt.Errorf("main: 0 not allowed; only positive integers"))
		}
	}

	app := tview.NewApplication().EnableMouse(true)
	root := tview.NewGrid()

	if len(durations) == 0 {
		s := widget.NewStopwatch()
		s.Changed = func() {
			app.Draw()
		}
		root.SetRows(0)
		root.AddItem(s, 0, 0, 1, 1, 0, 0, false)
		s.Start()
	} else {
		t := widget.NewTimer(durations[0])
		p := widget.NewProgressBar()
		t.SetVerticalAlign(widget.AlignDown)
		p.SetAlign(widget.AlignUp)
		t.Changed = func() {
			p.SetPercent(t.ElapsedSeconds() * 100 / t.TotalSeconds())
			app.Draw()
		}
		root.SetRows(0, 0)
		root.AddItem(t, 0, 0, 1, 1, 0, 0, false)
		root.AddItem(p, 1, 0, 1, 1, 0, 0, false)
		t.Start()
	}

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}

// ParseDuration returns the total number of seconds in dur, which must
// be of format [[hh:]mm:]ss.
func ParseDuration(dur string) (int, error) {
	var hr, min, sec int

	if m, err := regexp.MatchString(`^\d*$`, dur); m {
		if err != nil {
			return 0, err
		}
		sec, _ = strconv.Atoi(dur)
	} else if m, err := regexp.MatchString(`^\d+:\d{2}$`, dur); m {
		if err != nil {
			return 0, err
		}
		s := strings.Split(dur, ":")
		min, _ = strconv.Atoi(s[0])
		sec, _ = strconv.Atoi(s[1])
		// it's okay for minute field to be more than 60
		if err = checkField(sec, 0); err != nil {
			return 0, err
		}
	} else if m, err := regexp.MatchString(`^\d+:\d{2}:\d{2}$`, dur); m {
		if err != nil {
			return 0, err
		}
		s := strings.Split(dur, ":")
		hr, _ = strconv.Atoi(s[0])
		min, _ = strconv.Atoi(s[1])
		sec, _ = strconv.Atoi(s[2])
		if err = checkField(sec, min); err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("duration must be in [[hh:]mm:]ss format")
	}

	return (hr * 3600) + (min * 60) + sec, nil
}

// checkField returns error if sec/min field are not less than 60.
func checkField(sec, min int) error {
	var errmsg string
	if sec >= 60 {
		errmsg = "second's"
	}
	if min >= 60 && errmsg == "" {
		errmsg = "minute's"
	} else if min >= 60 {
		errmsg += " and minute's"
	}
	if errmsg != "" {
		return fmt.Errorf("%v field must be less than 60", errmsg)
	}
	return nil
}
