package task

import (
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ bubble.BubbleModel = &Model{}

type EditFinished struct {
	Content string
}

type Model struct {
	list.Model

	input   *textinput.Model
	editing *bool

	keymap *KeyMap
	repo   internal.Repo
	err    error
}

// IsLocked implements bubble.BubbleModel.
func (m *Model) IsLocked() bool {
	return *m.editing
}

// Error implements bubble.BubbleModel.
func (m *Model) Error() error {
	return m.err
}

func NewModel(k *KeyMap, r internal.Repo) *Model {
	ti := textinput.New()
	ti.Placeholder = "Task..."
	ti.CharLimit = 50
	ti.Width = 50
	ti.Prompt = ""

	editing := false
	l := list.New(
		[]list.Item{},
		CustomItemRender{
			Editor:  &ti,
			Editing: &editing,
		},
		50,
		10,
	)

	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(1, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)

	return &Model{
		keymap:  k,
		Model:   l,
		repo:    r,
		input:   &ti,
		editing: &editing,
	}
}

func (m *Model) View() string {

	if *m.editing {
		m.Title = "🖊️  Tasks"
	} else {
		m.Title = "📕 Tasks"
	}

	return m.Model.View()
}

func (m *Model) Fill(tasks ...entity.Task) {
	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}
	m.SetItems(items)
}

func (m *Model) Init() tea.Cmd {
	m.input.Focus()
	tasks, err := m.repo.List()
	if err != nil {
		return func() tea.Msg {
			return err
		}
	}
	m.Fill(tasks...)
	return nil
}

type editFinished struct{}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if *m.editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case msg.String() == "left" && m.input.Position() == 0,
				key.Matches(msg, m.keymap.EditExit):
				*m.editing = false
				return m.Update(editFinished{})
			}
		}

		var cmd tea.Cmd
		*m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {

	case EditFinished:
		if i, ok := m.SelectedItem().(entity.Task); ok {
			*m.editing = false
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
			if m.Index() != 0 {
				m.Select(m.Index() - 1)
			}

		case key.Matches(msg, m.keymap.Edit):
			if i, ok := m.SelectedItem().(entity.Task); ok {
				m.input.Focus()
				*m.editing = true
				m.input.SetValue(i.Description)
				m.input.SetCursor(0)
				return m, textinput.Blink
			}
		default:
			if len(msg.Runes) == 0 {
				break
			}

			if i, ok := m.SelectedItem().(entity.Task); ok {
				if i.Description == "" {
					*m.editing = true
					m.input.SetValue(string(msg.Runes))
					m.input.SetCursor(len(m.input.Value()))
					return m, textinput.Blink
				}
			}
		}
	case editFinished:
		if i, ok := m.SelectedItem().(entity.Task); ok {
			m.repo.Edit(&i, m.input.Value())
			m.Model.SetItem(m.Index(), i)
		}
	}

	return m, m.Init()
}
