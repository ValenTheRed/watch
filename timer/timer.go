package timer

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
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

// Timer component for Timer. Duration must never be zero.
type timer struct {
	*tview.TextView
	km                map[string]*help.Binding
	title             string
	duration, elapsed int
	running           bool
}

// newTimer returns a new timer.
func newTimer(duration int) *timer {
	t := &timer{
		TextView: tview.NewTextView(),
		title:    " Timer ",
		duration: duration,
		elapsed:  0,
		km: map[string]*help.Binding{
			"Reset": help.NewBinding(
				help.WithRune('r'), help.WithHelp("Reset"),
			),
			"Pause": help.NewBinding(
				help.WithRune('p'), help.WithHelp("Pause"),
			),
			"Start": help.NewBinding(
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
		SetTitle(t.title)

	return t
}

// Title returns the default title of t.
func (t *timer) Title() string {
	return t.title
}

func (t *timer) String() string {
	const (
		boundary = "┃"
		fill     = "█"
	)

	elapsed := utils.FormatSecond(t.elapsed)
	dur := utils.FormatSecond(t.duration)

	_, _, width, _ := t.GetInnerRect()
	// +2 is for the boundary chars at the either end of the progress
	// bar.
	width -= len(elapsed) + 2*tview.TabSize + 2 + 2*tview.TabSize + len(dur)
	percent := t.elapsed * 100 / t.duration
	fillLen := width * percent / 100

	return fmt.Sprintf(
		"\t%s\t%s\t%s\t", elapsed,
		strings.Join([]string{
			boundary,
			strings.Repeat(fill, fillLen),
			strings.Repeat(" ", width-fillLen),
			boundary,
		}, ""),
		dur,
	)
}

func (t *timer) IsTimeLeft() bool {
	return t.elapsed < t.duration
}

type interval struct {
	start, end time.Time
}

type Timer struct {
	*tview.Flex
	app *tview.Application

	Queue *queue
	Timer *timer

	stopMsg          chan struct{}
	pingMsg          chan interval
	timerSelectedMsg chan struct{}

	pingBuffer *beep.Buffer
}

// New returns a new Timer.
func New(durations []int, app *tview.Application) *Timer {
	// Init ping file
	// NOTE: error ignored
	streamer, format, _ := flac.Decode(bytes.NewReader(pingFile))
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	buf := beep.NewBuffer(format)
	buf.Append(streamer)

	// Init components
	timer, queue := newTimer(durations[0]), newQueue()
	for _, d := range durations {
		queue.addDuration(d)
	}

	t := &Timer{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(timer, 3, 0, true).
			AddItem(queue, 0, 1, false),
		app:   app,
		Timer: timer,
		// Queue will be used as the storage for all of the durations
		// information.
		Queue:      queue,
		pingBuffer: buf,
		// Channel is buffered because: `Stop()` -- which sends on
		// `stopMsg` -- will be called by the instance of `worker()`
		// started by `Start()`, which has it's `quit` channel
		// set to `stopMsg`; `Stop()` will block an unbuffered `stopMsg`.
		stopMsg:          make(chan struct{}, 1),
		pingMsg:          make(chan interval),
		timerSelectedMsg: make(chan struct{}),
	}

	t.Queue.setSelectFunc(func() {
		t.timerSelectedMsg <- struct{}{}
	})
	go t.queueControl()

	t.Timer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case t.Timer.km["Reset"].Rune():
			t.Reset()
			t.Timer.km["Start"].SetDisable(true)
			t.Timer.km["Pause"].SetDisable(false)
		case t.Timer.km["Pause"].Rune():
			if t.Timer.km["Pause"].IsEnabled() {
				t.Stop()
				t.Timer.km["Pause"].SetDisable(true)
				t.Timer.km["Start"].SetDisable(false)
			}
		case t.Timer.km["Start"].Rune():
			if t.Timer.km["Start"].IsEnabled() {
				t.Start()
				t.Timer.km["Start"].SetDisable(true)
				t.Timer.km["Pause"].SetDisable(false)
			}
		}
		return event
	})

	t.QueueTimerDraw()
	return t
}

// Keys returns the key bindings t uses for it's timer component.
func (t *Timer) Keys() []*help.Binding {
	return []*help.Binding{
		t.Timer.km["Start"],
		t.Timer.km["Pause"],
		t.Timer.km["Reset"],
	}
}

func (t *Timer) QueueTimerDraw() {
	go t.app.QueueUpdateDraw(func() {
		t.Timer.SetText(t.Timer.String())
	})
}

func (t *Timer) Start() {
	if t.Timer.running || !t.Timer.IsTimeLeft() {
		return
	}
	t.Timer.running = true
	go utils.Worker(func() {
		t.Timer.elapsed++
		t.QueueTimerDraw()
		if t.Timer.IsTimeLeft() {
			return
		}
		t.Stop()

		start := time.Now()
		speaker.Play(beep.Seq(
			t.pingBuffer.Streamer(0, t.pingBuffer.Len()),
			beep.Callback(func() {
				t.pingMsg <- interval{start, time.Now()}
			}),
		))
	}, t.stopMsg)
}

func (t *Timer) Stop() {
	if t.Timer.running {
		t.Timer.running = false
		t.stopMsg <- struct{}{}
	}
}

// queueControl is the handler that handles switching t
// - when user selects a timer from the queue
// - when current timer expires and the next timer from the queue needs
// to be played
// NOTE: queueControl needs to be run as a goroutine.
func (t *Timer) queueControl() {
	selectDone := time.Now()
	for {
		select {
		case interval := <-t.pingMsg:
			// If user selects a new timer within the time it takes for
			// ping sound to start and end, don't autostart next timer.
			if selectDone.Sub(interval.start) >= 0 &&
				interval.end.Sub(selectDone) >= 0 {
				break
			} else if err := t.Queue.queueNext(); err != nil {
				break
			}
			t.Timer.duration = t.Queue.getCurrentDuration()
			t.Reset()
		case <-t.timerSelectedMsg:
			t.Stop()
			t.Timer.duration = t.Queue.getCurrentDuration()
			t.Reset()
			selectDone = time.Now()
		}
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.Timer.elapsed = 0
	t.QueueTimerDraw()
	t.Start()
}
