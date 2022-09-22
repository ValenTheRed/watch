package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type Timer struct {
	*tview.Box

	// total is the total time in seconds.
	total int

	// elapsed is the time passed in seconds.
	elapsed int

	// Both determine the alignment of the timer text.
	verticalAlign, horizontalAlign int

	// stopCh will be used to signal to timer to stop ticking.
	stopCh chan struct{}

	// TextColor is the text color timer.
	TextColor tcell.Color

	// ShadowColor is the color for the shadow characters of the text.
	ShadowColor tcell.Color

	// done is an optional function that would be executed when timer
	// finishes.
	done func()
}

// NewTimer returns an initialised Timer that would counting down for
// duration seconds, and has it's text centered aligned both, vertically
// and horizontally.
func NewTimer(duration int) *Timer {
	return &Timer{
		Box:             tview.NewBox(),
		total:           duration,
		verticalAlign:   AlignCenter,
		horizontalAlign: tview.AlignCenter,
		// Timer will tick every second using Worker(). stopCh will be
		// used as the quit channel for Worker() and will be called from
		// the work() function being executed by Worker(). Since work()
		// is not executed in a go routine, signaling quit/stopCh
		// channel from work() would lead to a deadlock.
		stopCh:      make(chan struct{}, 1),
		TextColor:   tcell.ColorWhite,
		ShadowColor: tcell.ColorGrey,
	}
}

// SetHorizontalAlign sets the veritcal alignment of the text. Must be
// one of tview.AlignCenter, tview.AlignLeft or tview.AlignRight.
func (t *Timer) SetHorizontalAlign(align int) *Timer {
	t.horizontalAlign = align
	return t
}

// SetVerticalAlign sets the veritcal alignment of the text. Must be
// one of AlignCenter, AlignUp or AlignDown.
func (t *Timer) SetVerticalAlign(align int) *Timer {
	t.verticalAlign = align
	return t
}

// SetDoneFunc sets a handler which is called when the timer has
// finished.
func (t *Timer) SetDoneFunc(handler func()) *Timer {
	t.done = handler
	return t
}

// IsTimeLeft returns whether t has count down for duration it was set
// for.
func (t *Timer) IsTimeLeft() bool {
	return t.elapsed < t.total
}

// Start starts the Timer if time is left.
func (t *Timer) Start() *Timer {
	if !t.IsTimeLeft() {
		return t
	}
	go Worker(func() {
		t.elapsed++
		if t.IsTimeLeft() {
			return
		}
		t.Stop()
		if t.done != nil {
			t.done()
		}
	}, t.stopCh)
	return t
}

// Stop stops the Timer if time is left.
func (t *Timer) Stop() *Timer {
	if t.IsTimeLeft() {
		t.stopCh <- struct{}{}
	}
	return t
}

func (t *Timer) Draw(screen tcell.Screen) {
	t.DrawForSubclass(screen, t)

	text := SecondToANSIShadowWithLetters(t.total - t.elapsed)

	x, y, width, height := t.GetInnerRect()
	if t.verticalAlign == AlignCenter {
		y += getCenter(height, len(text))
	} else if t.verticalAlign == AlignDown {
		y += height - len(text)
	}
	if t.horizontalAlign == tview.AlignCenter {
		x += getCenter(width, runewidth.StringWidth(text[0]))
	} else if t.horizontalAlign == tview.AlignRight {
		x += width - runewidth.StringWidth(text[0])
	}

	shadowStyle := tcell.StyleDefault.Foreground(t.ShadowColor).Background(t.GetBackgroundColor())
	textStyle := tcell.StyleDefault.Foreground(t.TextColor).Background(t.GetBackgroundColor())

	for _, s := range text {
		i := 0
		for _, r := range s {
			style := textStyle
			if r != '█' {
				style = shadowStyle
			}
			screen.SetContent(x+i, y, r, nil, style)
			i++
		}
		y++
	}
}
