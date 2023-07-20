package task

import (
	"github.com/charmbracelet/bubbles/key"
)

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
func (m *Model) FullHelp() []key.Binding {
	if *m.editing {
		return []key.Binding{
			m.keymap.EditExit,
		}
	}
	return []key.Binding{
		m.keymap.Insert,
		m.keymap.Remove,
		m.keymap.Check,
		m.keymap.SwapUp,
		m.keymap.SwapDown,
		m.keymap.Up,
		m.keymap.Edit,
		m.keymap.Down,
	}
}

// ShortHelp implements help.KeyMap.
func (m *Model) ShortHelp() []key.Binding {
	if *m.editing {
		return []key.Binding{
			m.keymap.EditExit,
		}
	}

	return []key.Binding{
		m.keymap.Insert,
		m.keymap.Remove,
		m.keymap.Edit,
	}
}
