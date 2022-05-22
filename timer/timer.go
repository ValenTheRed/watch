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

type keyMap struct {
	Reset, Stop, Start *help.Binding
}

type Timer struct {
	*tview.TextView
	duration, timeLeft int
	running            bool
	stopMsg            chan struct{}
	title              string
	km                 keyMap

	app *tview.Application
}

func New(duration int, app *tview.Application) *Timer {
	t := &Timer{
		app: app,
		TextView: tview.NewTextView(),
		// Channel is buffered because: `Stop()` -- which sends on
		// `stopMsg` -- will be called by the instance of `worker()`
		// started by `Start()`, which has it's `quit` channel
		// set to `stopMsg`; `Stop()` will block an unbuffered `stopMsg`.
		stopMsg:  make(chan struct{}, 1),
		duration: duration,
		timeLeft: duration,
		title:    " Timer ",
		km: keyMap{
			Reset: help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			Stop: help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			Start: help.NewBinding(
				help.WithRune('s'),
				help.WithHelp("Start"),
				help.WithDisable(true),
			),
		},
	}

	t.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case t.km.Reset.Rune():
				t.Reset()
				t.km.Start.SetDisable(true)
				t.km.Stop.SetDisable(false)
			case t.km.Stop.Rune():
				if t.km.Stop.IsEnabled() {
					t.Stop()
					t.km.Stop.SetDisable(true)
					t.km.Start.SetDisable(false)
				}
			case t.km.Start.Rune():
				if t.km.Start.IsEnabled() {
					t.Start()
					t.km.Start.SetDisable(true)
					t.km.Stop.SetDisable(false)
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

func (t *Timer) Keys() []*help.Binding {
	return []*help.Binding{t.km.Reset, t.km.Start, t.km.Stop}
}

func (t *Timer) IsTimeLeft() bool {
	return t.timeLeft > 0
}

func (t *Timer) UpdateDisplay() {
	go t.app.QueueUpdateDraw(func() {
		t.SetText(utils.FormatSecond(t.timeLeft))
	})
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
				// go t.SetText(
				// 	fmt.Sprintf(
				// 		"Your %s's up!\n",
				// 		utils.FormatSecond(t.duration),
				// 	),
				// )
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