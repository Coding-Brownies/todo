package bubble

import (
	"strings"

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
	More     key.Binding
}

// FullHelp implements help.KeyMap.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{
		k.Quit,
		k.Check,
		k.Edit,
		k.EditExit,
		k.SwapUp,
		k.SwapDown,
		k.Insert,
		k.Remove,
		k.Up,
		k.Down,
	}}
}

// ShortHelp implements help.KeyMap.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Check,
		k.Edit,
		k.EditExit,
		k.Insert,
		k.Remove,
		k.More,
	}
}

func replaceSymbols(input string) string {
	unicodeMap := map[string]string{
		"ctrl":      "⌃",
		" ":         "␣",
		"up":        "↑",
		"down":      "↓",
		"left":      "←",
		"right":     "→",
		"shift":     "⇧",
		"tab":       "⇥",
		"backspace": "⌫",
		"delete":    "⌦",
		"enter":     "↵",
		"?":         "?",
	}
	words := strings.Split(input, "+")
	for i, word := range words {
		if unicodeValue, ok := unicodeMap[word]; ok {
			words[i] = unicodeValue
		}
	}
	return strings.Join(words, "+")
}

func NewKeyMap(cfg *config.Config) *KeyMap {
	return &KeyMap{
		Check: key.NewBinding(
			key.WithKeys(cfg.Check),
			key.WithHelp(replaceSymbols(cfg.Check), "(un)check "),
		),
		Quit: key.NewBinding(
			key.WithKeys(cfg.Quit),
			key.WithHelp(replaceSymbols(cfg.Quit), "quit "),
		),
		SwapUp: key.NewBinding(
			key.WithKeys(cfg.SwapUp),
			key.WithHelp(replaceSymbols(cfg.SwapUp), "swap up "),
		),
		SwapDown: key.NewBinding(
			key.WithKeys(cfg.SwapDown),
			key.WithHelp(replaceSymbols(cfg.SwapDown), "swap down "),
		),
		Remove: key.NewBinding(
			key.WithKeys(cfg.Remove),
			key.WithHelp(replaceSymbols(cfg.Remove), "remove "),
		),
		Insert: key.NewBinding(
			key.WithKeys(cfg.Insert),
			key.WithHelp(replaceSymbols(cfg.Insert), "insert "),
		),
		Edit: key.NewBinding(
			key.WithKeys(cfg.Edit),
			key.WithHelp(replaceSymbols(cfg.Edit), "edit "),
		),
		EditExit: key.NewBinding(
			key.WithKeys(cfg.EditExit),
			key.WithHelp(replaceSymbols(cfg.EditExit), "edit exit "),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("up ", ""),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("down ", ""),
		),
		More: key.NewBinding(
			key.WithKeys(cfg.More),
			key.WithHelp(replaceSymbols(cfg.More), "more "),
		),
	}
}
