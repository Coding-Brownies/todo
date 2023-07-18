package task

import (
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/bubble/components"
	"github.com/Coding-Brownies/todo/internal/bubble/components/snorkel"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubble.BubbleModel = &Model{}

type EditFinished struct {
	Content string
}

type Editor interface {
	Edit(string) tea.Cmd
}

type Model struct {
	list.Model

	editor  Editor
	editing bool

	keymap *KeyMap
	repo   internal.Repo
	err    error
}

func (m *Model) Map() help.KeyMap {
	return m.keymap
}

// Error implements bubble.BubbleModel.
func (m *Model) Error() error {
	return m.err
}

func NewModel(k *KeyMap, r internal.Repo, e Editor) *Model {
	l := list.New(
		[]list.Item{},
		components.CustomItemRender{},
		20,
		10,
	)

	l.Title = "📕 Tasks"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(1, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)

	return &Model{
		keymap: k,
		Model:  l,
		repo:   r,
		editor: e,
	}
}

func (m *Model) Fill(tasks ...entity.Task) {
	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}
	m.SetItems(items)
}

func (m *Model) Init() tea.Cmd {
	return func() tea.Msg {
		snorkel.Log("lolzone")

		tasks, err := m.repo.List()

		if err != nil {
			return err
		}
		m.Fill(tasks...)
		return nil
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editing {
		return m, nil
	}

	switch msg := msg.(type) {

	case EditFinished:
		if i, ok := m.SelectedItem().(entity.Task); ok {
			m.editing = false
			m.repo.Edit(&i, msg.Content)
		}

	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
		return m, nil

	case error:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch {
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

		case key.Matches(msg, m.keymap.Check):
			if i, ok := m.SelectedItem().(entity.Task); ok {
				m.repo.Check(&i)
				m.SetItem(m.Index(), i)
			}

		case key.Matches(msg, m.keymap.SwapUp):
			cur, ok := m.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.Select(m.Index() - 1)
			above, ok := m.SelectedItem().(entity.Task)
			if !ok {
				m.Select(m.Index() + 1)
				break
			}
			// Store changes synchronously
			m.repo.Swap(&cur, &above)

			m.SetItem(m.Index(), cur)
			m.SetItem(m.Index()+1, above)

		case key.Matches(msg, m.keymap.SwapDown):
			cur, ok := m.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.Select(m.Index() + 1)
			below, ok := m.SelectedItem().(entity.Task)
			if !ok {
				m.Select(m.Index() - 1)
				break
			}
			// Store changes synchronously
			m.repo.Swap(&cur, &below)

			m.SetItem(m.Index(), cur)
			m.SetItem(m.Index()-1, below)

		case key.Matches(msg, m.keymap.Insert):
			t := &entity.Task{}
			m.repo.Add(t)

			index := len(m.Items())
			m.InsertItem(index, *t)
			m.Select(index)

		case key.Matches(msg, m.keymap.Remove):
			// Store changes synchronously
			if i, ok := m.SelectedItem().(entity.Task); ok {
				m.repo.Delete(&i)
			}
			m.RemoveItem(m.Index())
			// if the index is out of bound set it back
			if m.Index() == len(m.Items()) {
				m.Select(m.Index() - 1)
			}

		case key.Matches(msg, m.keymap.Edit):
			if i, ok := m.SelectedItem().(entity.Task); ok {
				m.editing = true
				return m, tea.Batch(
					tea.HideCursor,
					m.editor.Edit(i.Description),
				)
			}
		}

	}

	return m, m.Init()
}
