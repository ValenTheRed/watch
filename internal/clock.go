package widget

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type Clock struct {
	*tview.Box

	// total is the total time in seconds.
	total int

	// elapsed is the time passed in seconds.
	elapsed int

	// Both determine the alignment of the clock text.
	verticalAlign, horizontalAlign int

	// stopCh will be used to signal to clock to stop ticking.
	stopCh chan struct{}

	// TextColor is the text color clock.
	TextColor tcell.Color

	// ShadowColor is the color for the shadow characters of the text.
	ShadowColor tcell.Color

	// done is an optional function that would be executed when clock
	// finishes.
	done func()

	// Format returns Clock value in ANSI Shadow font.
	Format func(second int) []string

	// value will be used by Format to generate the text of Clock to
	// draw.
	value func() int

	// Changed is an optional function that will be called when the
	// clock ticks. It is always safe to call app.Draw() from changed.
	Changed func()

	// Whether clock is running or not.
	running bool
}

// newClock returns a new Clock. It has horizontal and vertical aligment
// set to center, stopCh is uninitialised, value is the elapsed seconds,
// and Format is SecondToANSIShadowWithColons.
func newClock() *Clock {
	c := &Clock{
		Box:             tview.NewBox(),
		verticalAlign:   AlignCenter,
		horizontalAlign: tview.AlignCenter,
		TextColor:       tcell.ColorWhite,
		ShadowColor:     tcell.ColorGrey,
	}
	c.value = func() int {
		return c.elapsed
	}
	c.Format = SecondToANSIShadowWithColons
	return c
}

// NewTimer returns an initialised Clock that behaves like a timer. It
// counts down for duration seconds, and has it's text centered aligned
// both, vertically and horizontally. It uses
// SecondToANSIShadowWithLetters to format it's value.
func NewTimer(duration int) *Clock {
	c := newClock()
	c.total = duration
	// Timer will tick every second using Worker(). stopCh will be
	// used as the quit channel for Worker() and will be called from
	// the work() function being executed by Worker(). Since work()
	// is not executed in a go routine, signaling quit/stopCh
	// channel from work() would lead to a deadlock.
	c.stopCh = make(chan struct{}, 1)
	c.value = func() int {
		return c.total - c.elapsed
	}
	c.Format = SecondToANSIShadowWithLetters
	return c
}

// NewStopwatch returns an initialised Clock that behaves like a
// stopwatch. It has it's text centered aligned both, vertically and
// horizontally. It uses SecondToANSIShadowWithLetters to format it's
// value.
func NewStopwatch() *Clock {
	c := newClock()
	c.total = math.MaxInt
	// Stopwatch will never call Stop() from Worker(), so we don't need
	// stopCh to be buffered.
	c.stopCh = make(chan struct{})
	c.value = func() int {
		return c.elapsed
	}
	c.Format = SecondToANSIShadowWithLetters
	return c
}

// SetHorizontalAlign sets the veritcal alignment of the text. Must be
// one of tview.AlignCenter, tview.AlignLeft or tview.AlignRight.
func (c *Clock) SetHorizontalAlign(align int) *Clock {
	c.horizontalAlign = align
	return c
}

// SetVerticalAlign sets the veritcal alignment of the text. Must be
// one of AlignCenter, AlignUp or AlignDown.
func (c *Clock) SetVerticalAlign(align int) *Clock {
	c.verticalAlign = align
	return c
}

// SetDoneFunc sets a handler which is called when the clock has
// finished.
func (c *Clock) SetDoneFunc(handler func()) *Clock {
	c.done = handler
	return c
}

// IsTimeLeft returns whether c has count down for duration it was set
// for.
func (c *Clock) IsTimeLeft() bool {
	return c.elapsed < c.total
}

// SetElapsed sets the Clock's elapsed seconds to sec.
func (c *Clock) SetElapsed(sec int) *Clock {
	c.elapsed = sec
	if c.Changed != nil {
		go c.Changed()
	}
	return c
}

// Start starts the clock if time is left.
func (c *Clock) Start() *Clock {
	if c.running || !c.IsTimeLeft() {
		return c
	}
	c.running = true
	go Worker(func() {
		c.SetElapsed(c.elapsed + 1)
		if c.IsTimeLeft() {
			return
		}
		c.Stop()
		if c.done != nil {
			c.done()
		}
	}, c.stopCh)
	return c
}

// Stop signals clock to stop ticking.
func (c *Clock) Stop() *Clock {
	if c.running {
		c.running = false
		c.stopCh <- struct{}{}
	}
	return c
}

// Restart resets Clock's elapsed time to 0 and starts it again.
func (c *Clock) Restart() *Clock {
	c.Stop()
	c.SetElapsed(0)
	c.Start()
	return c
}

func (c *Clock) Draw(screen tcell.Screen) {
	c.DrawForSubclass(screen, c)

	text := c.Format(c.value())

	x, y, width, height := c.GetInnerRect()
	if c.verticalAlign == AlignCenter {
		y += getCenter(height, len(text))
	} else if c.verticalAlign == AlignDown {
		y += height - len(text)
	}
	if c.horizontalAlign == tview.AlignCenter {
		x += getCenter(width, runewidth.StringWidth(text[0]))
	} else if c.horizontalAlign == tview.AlignRight {
		x += width - runewidth.StringWidth(text[0])
	}

	shadowStyle := tcell.StyleDefault.Foreground(c.ShadowColor).Background(c.GetBackgroundColor())
	textStyle := tcell.StyleDefault.Foreground(c.TextColor).Background(c.GetBackgroundColor())

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
