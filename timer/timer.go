package timer

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

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

//go:embed "ping.flac"
var pingFile []byte

type keyMapTimer struct {
	Reset, Stop, Start *help.Binding
}

func (km keyMapTimer) Keys() []*help.Binding {
	return []*help.Binding{km.Reset, km.Start, km.Stop}
}

type Timer struct {
	*tview.TextView
	duration, timeLeft int
	running            bool
	stopMsg            chan struct{}
	title              string
	keyMap             keyMapTimer
}

func NewTimer(duration int) *Timer {
	t := &Timer{
		TextView: tview.NewTextView(),
		// Channel is buffered because: `Stop()` -- which sends on
		// `stopMsg` -- will be called by the instance of `worker()`
		// started by `Start()`, which has it's `quit` channel
		// set to `stopMsg`; `Stop()` will block an unbuffered `stopMsg`.
		stopMsg:  make(chan struct{}, 1),
		duration: duration,
		timeLeft: duration,
		title:    " Timer ",
		keyMap: keyMapTimer{
			Reset: help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			Stop: help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			Start: help.NewBinding(
				help.WithRune('s'), help.WithHelp("Start"),
			),
		},
	}
	t.keyMap.Start.SetDisable(true)

	t.
		SetChangedFunc(func() {
			wtc.app.Draw()
		}).
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetFocusFunc(focusFunc(t, t.keyMap)).
		SetBlurFunc(blurFunc(t)).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case t.keyMap.Reset.Rune():
				t.Reset()
				t.keyMap.Start.SetDisable(true)
				t.keyMap.Stop.SetDisable(false)
			case t.keyMap.Stop.Rune():
				if t.keyMap.Stop.IsEnabled() {
					t.Stop()
					t.keyMap.Stop.SetDisable(true)
					t.keyMap.Start.SetDisable(false)
				}
			case t.keyMap.Start.Rune():
				if t.keyMap.Start.IsEnabled() {
					t.Start()
					t.keyMap.Start.SetDisable(true)
					t.keyMap.Stop.SetDisable(false)
				}
			}
			return event
		}).
		SetTitle(t.title)

	t.UpdateDisplay()
	return t
}

func (t *Timer) Title() string {
	return t.title
}

func (t *Timer) IsTimeLeft() bool {
	return t.timeLeft > 0
}

func (t *Timer) UpdateDisplay() {
	t.SetText(utils.FormatSecond(t.timeLeft))
}

func (t *Timer) Start() {
	if !t.running && t.IsTimeLeft() {
		t.running = true
		go utils.Worker(func() {
			if t.IsTimeLeft() {
				t.timeLeft--
				t.UpdateDisplay()
			} else {
				t.Stop()
				// exec-ed in their own goroutine so that `stopMsg` can
				// get serviced before worker ticks.
				go t.SetText(
					fmt.Sprintf(
						"Your %s's up!\n",
						utils.FormatSecond(t.duration),
					),
				)
				go Ping(bytes.NewReader(pingFile))
			}
		}, t.stopMsg)
	}
}

func (t *Timer) Stop() {
	if t.running {
		t.running = false
		t.stopMsg <- struct{}{}
	}
}

func (t *Timer) Reset() {
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
