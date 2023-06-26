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
	Help     key.Binding
}

// FullHelp implements help.KeyMap.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Quit,
			k.Check,
			k.Insert,
			k.Remove,
		},
		{
			k.SwapUp,
			k.SwapDown,
			k.Up,
			k.Down,
		},
		{
			k.Edit,
			k.Help,
		},
	}
}

// ShortHelp implements help.KeyMap.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Check,
		k.Insert,
		k.Remove,
		k.Edit,
		k.Help,
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
			key.WithHelp(replaceSymbols(cfg.Check), "(un)check"),
		),
		Quit: key.NewBinding(
			key.WithKeys(cfg.Quit),
			key.WithHelp(replaceSymbols(cfg.Quit), "quit"),
		),
		SwapUp: key.NewBinding(
			key.WithKeys(cfg.SwapUp),
			key.WithHelp(replaceSymbols(cfg.SwapUp), "swap up"),
		),
		SwapDown: key.NewBinding(
			key.WithKeys(cfg.SwapDown),
			key.WithHelp(replaceSymbols(cfg.SwapDown), "swap down"),
		),
		Remove: key.NewBinding(
			key.WithKeys(cfg.Remove),
			key.WithHelp(replaceSymbols(cfg.Remove), "remove"),
		),
		Insert: key.NewBinding(
			key.WithKeys(cfg.Insert),
			key.WithHelp(replaceSymbols(cfg.Insert), "insert"),
		),
		Edit: key.NewBinding(
			key.WithKeys(cfg.Edit),
			key.WithHelp(replaceSymbols(cfg.Edit), "edit"),
		),
		EditExit: key.NewBinding(
			key.WithKeys(cfg.EditExit),
			key.WithHelp(replaceSymbols(cfg.EditExit), "to exit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp(replaceSymbols("up"), "go up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp(replaceSymbols("down"), "go down"),
		),
		Help: key.NewBinding(
			key.WithKeys(cfg.Help),
			key.WithHelp(replaceSymbols(cfg.Help), "toggle help"),
		),
	}
}
