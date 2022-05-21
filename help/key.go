package help

import "github.com/gdamore/tcell/v2"

// Design inspired from github.com/charmbracelet/bubbles/key

// Binding only supports ASCII printable characters as keybinds.
type Binding struct {
	key     tcell.Key
	char    rune
	disable bool
	help    string

	// Optional function that is called whenever disable is changed.
	handler func()
}

type BindingOpt func(*Binding)

// NewBinding returns a new Binding from a set of BindingOpt options.
func NewBinding(opts ...BindingOpt) *Binding {
	b := &Binding{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func WithRune(char rune) BindingOpt {
	return func(b *Binding) {
		b.char = char
		b.key = tcell.KeyRune
	}
}

func WithKey(key tcell.Key) BindingOpt {
	return func(b *Binding) {
		b.key = key
	}
}

func WithHelp(help string) BindingOpt {
	return func(b *Binding) {
		b.help = help
	}
}

func (b Binding) Rune() rune {
	return b.char
}

func (b Binding) Key() tcell.Key {
	return b.key
}

func (b Binding) Help() string {
	return b.help
}

func (b Binding) IsEnabled() bool {
	return !b.disable
}

func (b *Binding) SetDisable(opt bool) {
	b.disable = opt
	if b.handler != nil {
		b.handler()
	}
}

// SetDisableFunc sets handler as the function that is invoked whenever
// disable is changed.
func (b *Binding) SetDisableFunc(handler func()) {
	b.handler = handler
}

// DisableFunc returns the handler set with SetDisableFunc(), nil
// otherwise.
func (b *Binding) DisableFunc() func() {
	return b.handler
}
