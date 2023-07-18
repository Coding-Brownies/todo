package bin

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = &KeyMap{}

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Restore  key.Binding
	EmptyBin key.Binding
}

// FullHelp implements help.KeyMap.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Restore,
			k.EmptyBin,
		},
	}
}

// ShortHelp implements help.KeyMap.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Restore,
		k.EmptyBin,
	}
}
