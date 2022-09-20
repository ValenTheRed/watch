// A clock application with a stopwatch and a timer mode.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	usage = `usage: watch [-help] [duration]
A clock with a stopwatch and a timer.

Specify a duration to start a timer. Or, leave it alone to start a stopwatch.

optional arguments:
duration    supported formats - [[hh:]mm:]ss
-log        log to a file
-help	    display this help message and exit`

	// Global controller for the whole application.
	watch *Watch

	logArg bool
	debug  *log.Logger
)

func init() {
	flag.BoolVar(&logArg, "log", false, "log to a file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}

	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical
	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight

	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewHexColor(0x121014),
		ContrastBackgroundColor:     tcell.ColorBlue,
		MoreContrastBackgroundColor: tcell.ColorGreen,
		BorderColor:                 tcell.NewHexColor(0x51576a),
		TitleColor:                  tcell.NewHexColor(0x51576a),
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.NewHexColor(0xb3b3b3),
		SecondaryTextColor:          tcell.ColorYellow,
		TertiaryTextColor:           tcell.ColorGreen,
		InverseTextColor:            tcell.ColorBlue,
		ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
	}
}

func exitOnErr(err error) {
	if err != nil {
		debug.Fatalln(err)
	}
}

func main() {
	flag.Parse()

	var logFilename string
	if logArg {
		t := time.Now()
		logFilename = fmt.Sprintf(
			"watch_log_%d%02d%02d_%02d%02d%02d.log",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(),
		)
	} else {
		logFilename = os.DevNull
	}

	file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("main: %v\n", err)
	}
	defer file.Close()

	debug = log.New(file, "", log.LstdFlags|log.Lshortfile)

	durations := make([]int, len(flag.Args()))
	for i := range durations {
		durations[i], err = ParseDuration(flag.Arg(i))
		exitOnErr(err)
		if durations[i] == 0 {
			exitOnErr(fmt.Errorf("main: 0 not allowed; only positive integers"))
		}
	}

	watch = NewWatch(tview.NewApplication(), durations)

	if err := watch.Run(); err != nil {
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
