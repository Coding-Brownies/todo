package task

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

type Editor interface {
	Edit(string) (string, error)
}

type Model struct {
	list.Model

	editor Editor

	keymap KeyMap
	repo   internal.Repo
	Error  error
}

func NewModel(k KeyMap, r internal.Repo, e Editor) *Model {
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
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(0, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).Foreground(lipgloss.Color("#000000"))

	return &Model{
		keymap: k,
		Model:  l,
		repo:   r,
		editor: e,
	}
}

func (m Model) Fill(tasks ...entity.Task) {
	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}
	m.SetItems(items)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
		return m, nil

	case error:
		m.Error = msg
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

		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

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
				res, _ := m.editor.Edit(i.Description)

				// TODO add a check on max len
				m.repo.Edit(&i, res)
				m.SetItem(m.Index(), i)
			}

		}

	}
	return m, m.Init()
}
