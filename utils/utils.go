package utils

import (
	"fmt"
	"time"
)

func FormatSecond(s int) string {
	hrs := s / 3600
	min := (s / 60) % 60
	sec := s % 60
	return fmt.Sprintf("%02d:%02d:%02d", hrs, min, sec)
}

func worker(work func(), quit <-chan struct{}) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			work()
		case <-quit:
			return
		}
	}
}
