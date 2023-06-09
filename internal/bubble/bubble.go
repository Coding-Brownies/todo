package bubble

import (
	"fmt"

	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultWidth = 20

type model struct {
	keymap    KeyMap
	list      list.Model
	textInput textarea.Model
	err       error
	editing   bool
	bigHelp   bool
	listBin   bool
	repo      internal.Repo
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if key.Matches(msg, m.keymap.EditExit) {
				cur, ok := m.list.SelectedItem().(entity.Task)
				if !ok {
					break
				}
				// Store changes synchronously
				m.repo.Edit(&cur, m.textInput.Value())

				m.list.SetItem(m.list.Index(), cur)
				m.editing = false
			}
		case error:
			m.err = msg
			return m, nil
		}

		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	if m.listBin {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keymap.Bin):
				m.listBin = false // false o !m.listBin
				m.list.Title = DEFAULT_TITLE
				tasks, _ := m.repo.List()
				setList(tasks, &m.list)

			case key.Matches(msg, m.keymap.Restore):
				if i, ok := m.list.SelectedItem().(entity.Task); ok {
					m.repo.Restore(&i)
					tasks, _ := m.repo.ListBin()
					setList(tasks, &m.list)
				}

			case key.Matches(msg, m.keymap.Up):
				before := m.list.Index() - 1
				if before < 0 {
					before = 0
				}
				m.list.Select(before)

			case key.Matches(msg, m.keymap.Down):
				next := m.list.Index() + 1
				if next > len(m.list.Items())-1 {
					next = len(m.list.Items()) - 1
				}
				m.list.Select(next)

			case key.Matches(msg, m.keymap.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keymap.EmptyBin):
				m.repo.EmptyBin()
				setList([]entity.Task{}, &m.list)
			}

		case error:
			m.err = msg
			return m, nil
		}

		return m, m.Init()
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Bin):
			m.listBin = true
			m.list.Title = "🗑  Bin"
			tasks, _ := m.repo.ListBin()
			setList(tasks, &m.list)

		case key.Matches(msg, m.keymap.Help):
			m.bigHelp = !m.bigHelp

		case key.Matches(msg, m.keymap.Up):
			before := m.list.Index() - 1
			if before < 0 {
				before = 0
			}
			m.list.Select(before)

		case key.Matches(msg, m.keymap.Down):
			next := m.list.Index() + 1
			if next > len(m.list.Items())-1 {
				next = len(m.list.Items()) - 1
			}
			m.list.Select(next)

		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Check):
			if i, ok := m.list.SelectedItem().(entity.Task); ok {
				m.repo.Check(&i)
				m.list.SetItem(m.list.Index(), i)
			}

		case key.Matches(msg, m.keymap.SwapUp):
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.list.Select(m.list.Index() - 1)
			above, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				m.list.Select(m.list.Index() + 1)
				break
			}
			// Store changes synchronously
			m.repo.Swap(&cur, &above)

			m.list.SetItem(m.list.Index(), cur)
			m.list.SetItem(m.list.Index()+1, above)

		case key.Matches(msg, m.keymap.SwapDown):
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.list.Select(m.list.Index() + 1)
			below, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				m.list.Select(m.list.Index() - 1)
				break
			}
			// Store changes synchronously
			m.repo.Swap(&cur, &below)

			m.list.SetItem(m.list.Index(), cur)
			m.list.SetItem(m.list.Index()-1, below)

		case key.Matches(msg, m.keymap.Insert):
			t := &entity.Task{}
			m.repo.Add(t)

			index := len(m.list.Items())
			m.list.InsertItem(index, *t)
			m.list.Select(index)

		case key.Matches(msg, m.keymap.Remove):
			// Store changes synchronously
			if i, ok := m.list.SelectedItem().(entity.Task); ok {
				m.repo.Delete(&i)
			}
			m.list.RemoveItem(m.list.Index())
			// if the index is out of bound set it back
			if m.list.Index() == len(m.list.Items()) {
				m.list.Select(m.list.Index() - 1)
			}

		case key.Matches(msg, m.keymap.Edit):
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.editing = true
			m.textInput.SetValue(cur.Description)

		case key.Matches(msg, m.keymap.Undo):
			m.repo.Undo()
			tasks, _ := m.repo.List()
			setList(tasks, &m.list)

		default:
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			if cur.Description != "" {
				break
			}
			m.editing = true
			m.textInput.SetValue(msg.String())
		}
	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}
	return m, m.Init()
}

func (m model) View() string {
	if m.editing {
		return fmt.Sprintf(
			"\n✏️  Edit\n\n%s\n\n%s",
			m.textInput.View(),
			m.list.Help.ShortHelpView([]key.Binding{m.keymap.EditExit}),
		) + "\n"
	}
	if m.listBin {
		return fmt.Sprintf(
			"\n%s\n%s",
			m.list.View(),
			m.list.Help.ShortHelpView([]key.Binding{m.keymap.Quit, m.keymap.Bin, m.keymap.Restore, m.keymap.EmptyBin}),
		) + "\n"
	}
	help := m.list.Help.ShortHelpView(m.keymap.ShortHelp())
	if m.bigHelp {
		help = m.list.Help.FullHelpView(m.keymap.FullHelp())
	}
	return "\n" + m.list.View() + m.list.Styles.HelpStyle.Render(help)
}

const DEFAULT_TITLE = "📕 Tasks"

func New(cfg *config.Config, repo internal.Repo) *model {

	// build the input
	ta := textarea.New()
	ta.Placeholder = "Something todo..."
	ta.Focus()
	ta.CharLimit = 156
	ta.MaxHeight = 30

	l := list.New(
		[]list.Item{},
		customItemRender{},
		defaultWidth,
		10,
	)

	keyMap := NewKeyMap(cfg)

	// build the list
	l.Title = DEFAULT_TITLE
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(0, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2).Foreground(lipgloss.Color("#000000"))

	return &model{
		repo:      repo,
		list:      l,
		textInput: ta,
		keymap:    *keyMap,
	}
}

func (m *model) Run(tasks []entity.Task) error {

	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}
	m.list.SetItems(items)

	pg := tea.NewProgram(m)
	_, err := pg.Run()
	if err != nil {
		return err
	}
	return nil
}
