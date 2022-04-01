package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeSnapshot struct {
	TotalSeconds int
}

// Return error if sec/min field are not less than 60
func checkField(sec, min int, checkMin bool) error {
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
		if err = checkField(sec, min, false); err != nil {
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
		if err = checkField(sec, min, true); err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("duration must be in [[hh:]mm:]ss format")
	}

	return (hr * 3600) + (min * 60) + sec, nil
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
