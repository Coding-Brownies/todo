package bin

import (
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/bubble/components"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubble.BubbleModel = &Model{}

type Model struct {
	list.Model

	keymap *KeyMap
	repo   internal.Repo
	err    error
}

// IsLocked implements bubble.BubbleModel.
func (*Model) IsLocked() bool {
	return false
}

// Error implements bubble.BubbleModel.
func (m *Model) Error() error {
	return m.err
}

func NewModel(k *KeyMap, r internal.Repo) *Model {
	l := list.New(
		[]list.Item{},
		components.CustomItemRender{},
		20,
		10,
	)

	l.Title = "🗑  Bin"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(1, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)

	return &Model{
		keymap: k,
		Model:  l,
		repo:   r,
	}
}

func (m *Model) Fill(tasks ...entity.Task) {
	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}
	m.SetItems(items)
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	tasks, err := m.repo.ListBin()
	if err != nil {
		return func() tea.Msg {
			return err
		}
	}
	m.Fill(tasks...)
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
				m.Fill(tasks...)
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
			m.Fill()
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}
