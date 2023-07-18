package bubble

import (
	"strings"

	"github.com/Coding-Brownies/todo/config"
	"github.com/charmbracelet/bubbles/key"
	"golang.org/x/exp/slices"
)

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
	Undo     key.Binding
	Bin      key.Binding
	Restore  key.Binding
	EmptyBin key.Binding
}

func replaceSymbols(inputs []string) string {
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
		Check: key.NewBinding(
			WithKeys(cfg.Check...),
			key.WithHelp(replaceSymbols(cfg.Check), "(un)check"),
		),
		Quit: key.NewBinding(
			WithKeys(cfg.Quit...),
			key.WithHelp(replaceSymbols(cfg.Quit), "quit"),
		),
		SwapUp: key.NewBinding(
			WithKeys(cfg.SwapUp...),
			key.WithHelp(replaceSymbols(cfg.SwapUp), "swap up"),
		),
		SwapDown: key.NewBinding(
			WithKeys(cfg.SwapDown...),
			key.WithHelp(replaceSymbols(cfg.SwapDown), "swap down"),
		),
		Remove: key.NewBinding(
			WithKeys(cfg.Remove...),
			key.WithHelp(replaceSymbols(cfg.Remove), "remove"),
		),
		Insert: key.NewBinding(
			WithKeys(cfg.Insert...),
			key.WithHelp(replaceSymbols(cfg.Insert), "insert"),
		),
		Edit: key.NewBinding(
			WithKeys(cfg.Edit...),
			key.WithHelp(replaceSymbols(cfg.Edit), "edit"),
		),
		EditExit: key.NewBinding(
			WithKeys(cfg.EditExit...),
			key.WithHelp(replaceSymbols(cfg.EditExit), "to exit"),
		),
		Up: key.NewBinding(
			WithKeys(cfg.Up...),
			key.WithHelp(replaceSymbols(cfg.Up), "go up"),
		),
		Down: key.NewBinding(
			WithKeys(cfg.Down...),
			key.WithHelp(replaceSymbols(cfg.Down), "go down"),
		),
		Help: key.NewBinding(
			WithKeys(cfg.Help...),
			key.WithHelp(replaceSymbols(cfg.Help), "toggle help"),
		),
		Undo: key.NewBinding(
			WithKeys(cfg.Undo...),
			key.WithHelp(replaceSymbols(cfg.Undo), "undo"),
		),
		Bin: key.NewBinding(
			WithKeys(cfg.Cycle...),
			key.WithHelp(replaceSymbols(cfg.Cycle), "toggle bin"),
		),
		Restore: key.NewBinding(
			WithKeys(cfg.Restore...),
			key.WithHelp(replaceSymbols(cfg.Restore), "restore"),
		),
		EmptyBin: key.NewBinding(
			WithKeys(cfg.EmptyBin...),
			key.WithHelp(replaceSymbols(cfg.EmptyBin), "empty the bin"),
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
