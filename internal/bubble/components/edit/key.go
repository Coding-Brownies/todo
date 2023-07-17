package edit

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Exit key.Binding
}

// FullHelp implements help.KeyMap.
func (*KeyMap) FullHelp() [][]key.Binding {
	panic("unimplemented")
}

// ShortHelp implements help.KeyMap.
func (*KeyMap) ShortHelp() []key.Binding {
	panic("unimplemented")
}
