package bin

import (
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble/components"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &Model{}

type Model struct {
	list.Model

	keymap KeyMap
	repo   internal.Repo
	Error  error
}

func NewModel(k KeyMap, r internal.Repo) *Model {
	l := list.New(
		[]list.Item{},
		components.CustomItemRender{},
		20,
		10,
	)

	l.Title = "🗑 Bin"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(0, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).Foreground(lipgloss.Color("#000000"))

	return &Model{
		keymap: k,
		Model:  l,
		repo:   r,
	}
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keymap.Restore):
			if i, ok := m.SelectedItem().(entity.Task); ok {
				m.repo.Restore(&i)
				tasks, _ := m.repo.ListBin()
				components.SetList(tasks, &m.Model)
			}

		case key.Matches(msg, m.keymap.Up):
			before := m.Index() - 1
			if before < 0 {
				before = 0
			}
			m.Select(before)

		case key.Matches(msg, m.keymap.Down):
			next := m.Index() + 1
			if next > len(m.Items())-1 {
				next = len(m.Items()) - 1
			}
			m.Select(next)

		case key.Matches(msg, m.keymap.EmptyBin):
			m.repo.EmptyBin()
			components.SetList([]entity.Task{}, &m.Model)
		}

	case error:
		m.Error = msg
		return m, nil
	}

	return m, m.Init()
}
