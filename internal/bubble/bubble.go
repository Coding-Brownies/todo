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
	repo      internal.Repo
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		// case in which the m.bigHelp field of the model is changed
		case key.Matches(msg, m.keymap.Help):
			if m.editing {
				break
			}
			m.bigHelp = !m.bigHelp

		case key.Matches(msg, m.keymap.Up):
			if m.editing {
				break
			}
			before := m.list.Index() - 1
			if before < 0 {
				before = 0
			}
			m.list.Select(before)

		case key.Matches(msg, m.keymap.Down):
			if m.editing {
				break
			}
			next := m.list.Index() + 1
			if next > len(m.list.Items())-1 {
				next = len(m.list.Items()) - 1
			}
			m.list.Select(next)

		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Check):
			if m.editing {
				break
			}
			if i, ok := m.list.SelectedItem().(entity.Task); ok {
				//Store changes synchronously
				m.repo.Check(i.ID)

				i.Done = !i.Done
				m.list.SetItem(m.list.Index(), i)
			}

		case key.Matches(msg, m.keymap.SwapUp):
			if m.editing {
				break
			}
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
			m.repo.Swap(cur.ID, above.ID)

			m.list.SetItem(m.list.Index(), cur)
			m.list.SetItem(m.list.Index()+1, above)

		case key.Matches(msg, m.keymap.SwapDown):
			if m.editing {
				break
			}
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
			m.repo.Swap(cur.ID, below.ID)

			m.list.SetItem(m.list.Index(), cur)
			m.list.SetItem(m.list.Index()-1, below)

		case key.Matches(msg, m.keymap.Insert):
			if m.editing {
				break
			}
			// Store changes synchronously
			m.repo.Add(&entity.Task{})

			m.list.InsertItem(len(m.list.Items())-1, entity.Task{})
			m.list.Select(len(m.list.Items())-1)

		case key.Matches(msg, m.keymap.Remove):
			if m.editing {
				break
			}
			// Store changes synchronously
			if i, ok := m.list.SelectedItem().(entity.Task); ok {
				m.repo.Delete(i.ID)
			}

			if m.list.Index() == len(m.list.Items())-1 && len(m.list.Items()) > 1 {
				m.list.RemoveItem(m.list.Index())
				m.list.Select(m.list.Index() - 1)
			} else {
				m.list.RemoveItem(m.list.Index())
			}

		case key.Matches(msg, m.keymap.Edit):
			if m.editing {
				break
			}
			m.editing = true
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.textInput.SetValue(cur.Description)

		case key.Matches(msg, m.keymap.EditExit):
			if !m.editing {
				break
			}
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			// Store changes synchronously
			m.repo.Edit(cur.ID, m.textInput.Value())

			m.list.SetItem(m.list.Index(), entity.Task{
				Done:        cur.Done,
				Description: m.textInput.Value(),
			})
			m.editing = false
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd

	if m.editing {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, m.Init()
}

func (m model) View() string {
	if m.editing {
		return fmt.Sprintf(
			"\n%s\n\n%s",
			m.textInput.View(),
			m.list.Help.ShortHelpView([]key.Binding{m.keymap.EditExit}),
		) + "\n"
	}

	help := m.list.Help.ShortHelpView(m.keymap.ShortHelp())
	if m.bigHelp {
		help = m.list.Help.FullHelpView(m.keymap.FullHelp())
	}

	return m.list.View() +
		m.list.Styles.HelpStyle.Render(help)
}

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
	l.Title = ""
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
