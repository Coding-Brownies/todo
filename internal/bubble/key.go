package bubble

import (
	"github.com/Coding-Brownies/todo/config"
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
	Edit     key.Binding
	EditExit key.Binding
	Up       key.Binding
	Down     key.Binding
}

// FullHelp implements help.KeyMap.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.Edit,
	}}
}

// ShortHelp implements help.KeyMap.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Edit,
	}
}

func NewKeyMap(cfg *config.Config) *KeyMap {
	return &KeyMap{
		Check: key.NewBinding(
			key.WithKeys(cfg.Check),
			key.WithHelp(cfg.Check, "(un)check the tasks"),
		),
		Quit: key.NewBinding(
			key.WithKeys(cfg.Quit),
			key.WithHelp(cfg.Quit, "quit"),
		),
		SwapUp: key.NewBinding(
			key.WithKeys(cfg.SwapUp),
			key.WithHelp(cfg.SwapUp, "swap up"),
		),
		SwapDown: key.NewBinding(
			key.WithKeys(cfg.SwapDown),
			key.WithHelp(cfg.SwapDown, "swap down"),
		),
		Remove: key.NewBinding(
			key.WithKeys(cfg.Remove),
			key.WithHelp(cfg.Remove, "remove"),
		),
		Insert: key.NewBinding(
			key.WithKeys(cfg.Insert),
			key.WithHelp(cfg.Insert, "insert a new task"),
		),
		Edit: key.NewBinding(
			key.WithKeys(cfg.Edit),
			key.WithHelp(cfg.Edit, "edit"),
		),
		EditExit: key.NewBinding(
			key.WithKeys(cfg.EditExit),
			key.WithHelp(cfg.EditExit, "edit exit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("up", ""),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("down", ""),
		),
	}
}
