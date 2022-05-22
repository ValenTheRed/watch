package help

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeyMaper interface {
	// Keys will return the key Binding in the order that the widgets
	// want them to be displayed.
	Keys() []*Binding
}

type keyMap struct{}

type HelpView struct {
	*tview.TextView
	globals, locals []*Binding
	title           string
	km              keyMap

	app *tview.Application
}

func NewHelpView(app *tview.Application) *HelpView {
	hv := &HelpView{
		app: app,
		TextView: tview.NewTextView(),
		title:    " Help ",
		km:       keyMap{},
	}
	hv.
		SetTextAlign(tview.AlignCenter).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault).
		SetTitle(hv.title)

	return hv
}

func (hv *HelpView) Title() string {
	return hv.title
}

func (hv *HelpView) Keys() []*Binding {
	return []*Binding{}
}

func (hv *HelpView) UpdateDisplay() {
	sep := " • "

	var view []string
	for _, bindings := range [][]*Binding{
		hv.locals,
		hv.globals,
	} {
		for _, b := range bindings {
			if !b.IsEnabled() {
				continue
			}
			var key string
			if b.Key() == tcell.KeyRune {
				key = string(b.Rune())
			} else {
				key = tcell.KeyNames[b.Key()]
			}
			view = append(view, key+sep+b.Help())
		}
	}

	go hv.app.QueueUpdateDraw(func() {
		hv.SetText(strings.Join(view, "\t"))
	})
}

func (hv *HelpView) SetGlobals(km KeyMaper) {
	hv.setBindings(&hv.globals, km)
}

func (hv *HelpView) UnsetGlobals() {
	hv.unsetBindings(&hv.globals)
}

func (hv *HelpView) SetLocals(km KeyMaper) {
	hv.setBindings(&hv.locals, km)
}

func (hv *HelpView) UnsetLocals() {
	hv.unsetBindings(&hv.locals)
}

func (hv *HelpView) setBindings(bindings *[]*Binding, km KeyMaper) {
	*bindings = km.Keys()
	for _, binding := range *bindings {
		binding.SetDisableFunc(hv.UpdateDisplay)
	}
}

func (hv *HelpView) unsetBindings(bindings *[]*Binding) {
	for _, binding := range *bindings {
		binding.SetDisableFunc(nil)
	}
	*bindings = nil
}