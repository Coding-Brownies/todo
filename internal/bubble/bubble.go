package bubble

import (
	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &model{}

type BubbleModel interface {
	tea.Model
	Map() help.KeyMap
	Error() error
}

type editFinished struct{}

type model struct {
	repo   internal.Repo
	keymap *KeyMap
	err    error

	help help.Model

	// components
	cur    int
	models []BubbleModel
}

func (m model) Init() tea.Cmd {
	return m.models[m.cur].Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case editFinished:
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Bin):
			m.cur = (m.cur + 1) % len(m.models)
			return m, m.Init()

		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	_, cmd := m.models[m.cur].Update(msg)

	return m, cmd
}

func (m model) View() string {
	cur := m.models[m.cur]
	view := cur.View()

	h := m.help.View(cur.Map())

	return view + "\n" +
		list.DefaultStyles().
			HelpStyle.PaddingLeft(2).
			Foreground(lipgloss.Color("#000000")).Render(h)
}

func New(cfg *config.Config, repo internal.Repo, keyMap *KeyMap, models ...BubbleModel) *model {
	m := &model{
		keymap: keyMap,
		repo:   repo,
		models: models,
		help:   help.New(),
	}
	return m
}

func (m *model) Run() error {
	pg := tea.NewProgram(m)
	_, err := pg.Run()
	if err != nil {
		return err
	}
	return nil
}
