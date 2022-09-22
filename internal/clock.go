package widget

import (
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
}

// newClock returns a Clock with horizontal and vertical aligment set to
// center, and an uninitialised stopCh.
func newClock() *Clock {
	return &Clock{
		Box:             tview.NewBox(),
		verticalAlign:   AlignCenter,
		horizontalAlign: tview.AlignCenter,
		TextColor:       tcell.ColorWhite,
		ShadowColor:     tcell.ColorGrey,
	}
}

// NewTimer returns an initialised Clock that behaves like a timer. It
// counts down for duration seconds, and has it's text centered aligned
// both, vertically and horizontally.
func NewTimer(duration int) *Clock {
	c := newClock()
	c.total = duration
	// Timer will tick every second using Worker(). stopCh will be
	// used as the quit channel for Worker() and will be called from
	// the work() function being executed by Worker(). Since work()
	// is not executed in a go routine, signaling quit/stopCh
	// channel from work() would lead to a deadlock.
	c.stopCh = make(chan struct{}, 1)
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

// Start starts the clock if time is left.
func (c *Clock) Start() *Clock {
	if !c.IsTimeLeft() {
		return c
	}
	go Worker(func() {
		c.elapsed++
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

// Stop stops the clock if time is left.
func (c *Clock) Stop() *Clock {
	if c.IsTimeLeft() {
		c.stopCh <- struct{}{}
	}
	return c
}

func (c *Clock) Draw(screen tcell.Screen) {
	c.DrawForSubclass(screen, c)

	text := SecondToANSIShadowWithLetters(c.total - c.elapsed)

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
