package task

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Check    key.Binding
	Quit     key.Binding
	SwapUp   key.Binding
	SwapDown key.Binding
	Remove   key.Binding
	Insert   key.Binding
	Up       key.Binding
	Down     key.Binding
	Edit     key.Binding
}

// FullHelp implements help.KeyMap.
func (*KeyMap) FullHelp() [][]key.Binding {
	panic("unimplemented")
}

// ShortHelp implements help.KeyMap.
func (*KeyMap) ShortHelp() []key.Binding {
	panic("unimplemented")
}
