package widget

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeyMap struct {
	Key, Desc string
}

type HelpView struct {
	*tview.TextView

	// keys is a keymap.
	keys []KeyMap

	// Styles for the a key's shortcut, it's help description, and the
	// separator character between two keymap.
	keyStyle, descStyle, separatorStyle tcell.Style
}

// parseStyleToTag converts style to tview color tags.
func parseStyleToTag(style tcell.Style) string {
	var s strings.Builder
	fg, bg, attr := style.Decompose()
	s.WriteRune('[')
	// The guard is necessary since the value of Hex() is -1 even for
	// tcell.ColorDefault, and -1 cannot be redered by tview. This
	// forces the user to define some value for foreground/background.
	// This would be alright, since foreground could be
	// tview.Styles.PrimaryTextColor and background could be
	// tview.Styles.PrimitiveBackgroundColor but (on my terminal at
	// least) but since tview.Styles.PrimitiveBackgroundColor is
	// tcell.ColorBlack, TextView screen actually rendering the black
	// color as opposed to the default color that's usually used for the
	// widgets.
	if hex := fg.Hex(); hex != -1 {
		s.WriteString(fmt.Sprintf("#%06x", hex))
	}
	s.WriteRune(':')
	if hex := bg.Hex(); hex != -1 {
		s.WriteString(fmt.Sprintf("#%06x", hex))
	}
	s.WriteRune(':')
	var AttrIs = func(mask tcell.AttrMask) bool {
		return attr&mask != 0
	}
	if !AttrIs(tcell.AttrNone) || !AttrIs(tcell.AttrInvalid) {
		if AttrIs(tcell.AttrBold) {
			s.WriteRune('b')
		}
		if AttrIs(tcell.AttrBlink) {
			s.WriteRune('l')
		}
		if AttrIs(tcell.AttrItalic) {
			s.WriteRune('i')
		}
		if AttrIs(tcell.AttrDim) {
			s.WriteRune('d')
		}
		if AttrIs(tcell.AttrReverse) {
			s.WriteRune('r')
		}
		if AttrIs(tcell.AttrUnderline) {
			s.WriteRune('u')
		}
		if AttrIs(tcell.AttrStrikeThrough) {
			s.WriteRune('s')
		}
	}
	s.WriteRune(']')
	return s.String()
}
