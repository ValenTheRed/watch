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

// NewHelpView returns a new HelpView.
func NewHelpView(keymaps []KeyMap) *HelpView {
	hv := &HelpView{keys: keymaps}
	hv.TextView = tview.NewTextView().SetText(hv.toTextViewString())
	return hv
}

// SetKeyStyle sets the style of keymap's key shortcut.
// NOTE: key style will inherit the background of HelpView, so no need to
// specify that.
func (hv *HelpView) SetKeyStyle(style tcell.Style) *HelpView {
	hv.keyStyle = style
	hv.TextView.SetText(hv.toTextViewString())
	return hv
}

// SetDescStyle sets the style of keymap's key description. NOTE: desc
// style will inherit the background of HelpView, so no need to
// specify that.
func (hv *HelpView) SetDescStyle(style tcell.Style) *HelpView {
	hv.descStyle = style
	hv.TextView.SetText(hv.toTextViewString())
	return hv
}

// SetSeparatorStyle sets the style of separator between two keymap.
// NOTE: separator style will inherit the background of HelpView, so no
// need to specify that.
func (hv *HelpView) SetSeparatorStyle(style tcell.Style) *HelpView {
	hv.separatorStyle = style
	hv.TextView.SetText(hv.toTextViewString())
	return hv
}

// toTextViewString converts the keymaps in hv to string that can be
// processed by the embedded tview.TextView.
func (hv *HelpView) toTextViewString() string {
	const separator = '•'

	keyTag := parseStyleToTag(hv.keyStyle)
	descTag := parseStyleToTag(hv.descStyle)
	separatorTag := parseStyleToTag(hv.separatorStyle)

	var s strings.Builder
	for i, km := range hv.keys {
		s.WriteString(keyTag)
		s.WriteString(km.Key)
		s.WriteString("[-:-:-] ")
		s.WriteString(descTag)
		s.WriteString(km.Desc)
		// write separator only if this isn't the last keybind.
		if i < len(hv.keys)-1 {
			s.WriteString("[-:-:-] ")
			s.WriteString(separatorTag)
			s.WriteRune(separator)
			s.WriteString("[-:-:-] ")
		}
	}
	return s.String()
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
