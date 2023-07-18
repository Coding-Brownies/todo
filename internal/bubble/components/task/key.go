package task

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Check    key.Binding
	SwapUp   key.Binding
	SwapDown key.Binding
	Remove   key.Binding
	Insert   key.Binding
	Up       key.Binding
	Down     key.Binding
	Edit     key.Binding
	EditExit key.Binding
}

// FullHelp implements help.KeyMap.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Insert,
			k.Remove,
			k.Check,
			k.SwapUp,
		},
		{
			k.SwapDown,
			k.Up,
			k.Edit,
			k.Down,
		},
		{
			k.EditExit,
		},
	}
}

// ShortHelp implements help.KeyMap.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Insert,
		k.Remove,
		k.Up,
		k.Edit,
		k.Down,
	}
}
