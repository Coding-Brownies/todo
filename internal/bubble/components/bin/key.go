package bin

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Restore  key.Binding
	EmptyBin key.Binding
}

// FullHelp implements help.KeyMap.
func (m *Model) FullHelp() []key.Binding {
	return []key.Binding{
		m.keymap.Up,
		m.keymap.Down,
		m.keymap.Restore,
		m.keymap.EmptyBin,
	}
}

// ShortHelp implements help.KeyMap.
func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keymap.Restore,
		m.keymap.EmptyBin,
	}
}
