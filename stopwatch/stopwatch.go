package stopwatch

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ValenTheRed/watch/help"
	"github.com/ValenTheRed/watch/utils"
)

type keyMap struct {
	Reset, Stop, Start *help.Binding
}

type Stopwatch struct {
	*tview.TextView
	elapsed int
	running bool
	stopMsg chan struct{}
	title   string
	km      keyMap
}

func New() *Stopwatch {
	sw := &Stopwatch{
		TextView: tview.NewTextView(),
		stopMsg:  make(chan struct{}),
		title:    " Stopwatch ",
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

	sw.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case sw.km.Reset.Rune():
				sw.Reset()
				sw.km.Start.SetDisable(true)
				sw.km.Stop.SetDisable(false)
			case sw.km.Stop.Rune():
				if sw.km.Stop.IsEnabled() {
					sw.Stop()
					sw.km.Stop.SetDisable(true)
					sw.km.Start.SetDisable(false)
				}
			case sw.km.Start.Rune():
				if sw.km.Start.IsEnabled() {
					sw.Start()
					sw.km.Start.SetDisable(true)
					sw.km.Stop.SetDisable(false)
				}
			}
			return event
		}).
		SetTitle(sw.title)

	sw.UpdateDisplay()
	return sw
}

func (sw *Stopwatch) Title() string {
	return sw.title
}

func (sw *Stopwatch) Keys() []*help.Binding {
	return []*help.Binding{sw.km.Reset, sw.km.Start, sw.km.Stop}
}

func (sw *Stopwatch) UpdateDisplay() {
	sw.SetText(utils.FormatSecond(sw.elapsed))
}

func (sw *Stopwatch) Start() {
	if !sw.running {
		sw.running = true
		go utils.Worker(func() {
			sw.elapsed++
			sw.UpdateDisplay()
		}, sw.stopMsg)
	}
}

func (sw *Stopwatch) Stop() {
	if sw.running {
		sw.running = false
		sw.stopMsg <- struct{}{}
	}
}

func (sw *Stopwatch) Reset() {
	sw.Stop()
	sw.elapsed = 0
	sw.UpdateDisplay()
	sw.Start()
}
