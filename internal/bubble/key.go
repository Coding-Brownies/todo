package bubble

import (
	"strings"

	"github.com/Coding-Brownies/todo/config"
	"github.com/charmbracelet/bubbles/key"
	"golang.org/x/exp/slices"
)

type KeyMap struct {
	Quit  key.Binding
	Help  key.Binding
	Undo  key.Binding
	Cycle key.Binding
}

func ReplaceSymbols(inputs []string) string {
	unicodeMap := map[string]string{
		"ctrl":      "⌃",
		"space":     "␣",
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
	res := []string{}
	for _, input := range inputs {
		words := strings.Split(input, "+")
		for i, word := range words {
			if unicodeValue, ok := unicodeMap[word]; ok {
				words[i] = unicodeValue
			}
		}
		res = append(res, strings.Join(words, "+"))
	}
	return strings.Join(res, "/")
}

func NewKeyMap(cfg *config.Config) *KeyMap {
	return &KeyMap{
		Quit: key.NewBinding(
			WithKeys(cfg.Quit...),
			key.WithHelp(ReplaceSymbols(cfg.Quit), "quit"),
		),
		Help: key.NewBinding(
			WithKeys(cfg.Help...),
			key.WithHelp(ReplaceSymbols(cfg.Help), "toggle help"),
		),
		Undo: key.NewBinding(
			WithKeys(cfg.Undo...),
			key.WithHelp(ReplaceSymbols(cfg.Undo), "undo"),
		),
		Cycle: key.NewBinding(
			WithKeys(cfg.Cycle...),
			key.WithHelp(ReplaceSymbols(cfg.Cycle), "toggle views"),
		),
	}
}

func WithKeys(keys ...string) key.BindingOpt {
	return func(b *key.Binding) {
		if i := slices.Index(keys, "space"); i != -1 {
			keys[i] = " "
		}
		b.SetKeys(keys...)
	}
}

// FullHelp implements help.KeyMap.
func (m *model) FullHelp() []key.Binding {
	return []key.Binding{
		m.keymap.Cycle,
		m.keymap.Help,
		m.keymap.Quit,
		m.keymap.Undo,
	}
}

// ShortHelp implements help.KeyMap.
func (m *model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keymap.Quit,
		m.keymap.Help,
	}
}

// DevideIntoColumns divides the input slice into multiple columns of the given number of rows each.
func DevideIntoColumns(bindings []key.Binding, rows int) [][]key.Binding {
	totalItems := len(bindings)

	columns := (totalItems + rows - 1) / rows // Round up the division
	result := make([][]key.Binding, columns)

	for i := 0; i < columns; i++ {
		start := i * rows
		end := (i + 1) * rows
		if end > totalItems {
			end = totalItems
		}
		result[i] = bindings[start:end]
	}

	return result
}
