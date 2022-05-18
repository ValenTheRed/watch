package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:embed "ping.flac"
var pingFile []byte

type timer struct {
	*tview.TextView
	duration, timeLeft int
	running            bool
	stopMsg            chan struct{}
}

func NewTimer(duration int) *timer {
	t := &timer{
		TextView: tview.NewTextView(),
		// Channel is buffered because: `Stop()` -- which sends on
		// `stopMsg` -- will be called by the instance of `worker()`
		// started by `Start()`, which has it's `quit` channel
		// set to `stopMsg`; `Stop()` will block an unbuffered `stopMsg`.
		stopMsg:  make(chan struct{}, 1),
		duration: duration,
		timeLeft: duration,
	}
	t.TextView.
		SetTextAlign(tview.AlignCenter).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'r':
				t.Reset()
			case 'p':
				t.Stop()
			case 's':
				t.Start()
			}
			return event
		}).
		SetTitle("Timer").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)

	return t
}

func (t *timer) UpdateDisplay() {
	t.SetText(FormatSecond(t.timeLeft))
}

func (t *timer) Start() {
	if !t.running {
		t.running = true
		go worker(func() {
			if t.timeLeft > 0 {
				t.timeLeft--
				t.UpdateDisplay()
			} else {
				t.Stop()
				// exec-ed in their own goroutine so that `stopMsg` can
				// get serviced before worker ticks.
				go t.SetText(
					fmt.Sprintf("Your %s's up!\n", FormatSecond(t.duration)),
				)
				go Ping(bytes.NewReader(pingFile))
			}
		}, t.stopMsg)
	}
}

func (t *timer) Stop() {
	if t.running {
		t.running = false
		t.stopMsg <- struct{}{}
	}
}

func (t *timer) Reset() {
	t.Stop()
	t.timeLeft = t.duration
	t.UpdateDisplay()
	t.Start()
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
