package bubble

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultWidth = 20

type keymap struct {
	check    key.Binding
	quit     key.Binding
	swapUp   key.Binding
	swapDown key.Binding
	remove   key.Binding
	insert   key.Binding
	edit     key.Binding
	editExit key.Binding
	up       key.Binding
	down     key.Binding
}

type model struct {
	keymap    keymap
	list      list.Model
	textInput textarea.Model
	err       error
	editing   bool
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
		case key.Matches(msg, m.keymap.up):
			before := m.list.Index() - 1
			if before < 0 {
				before = 0
			}
			m.list.Select(before)

		case key.Matches(msg, m.keymap.down):
			next := m.list.Index() + 1
			if next > len(m.list.Items())-1 {
				next = len(m.list.Items()) - 1
			}
			m.list.Select(next)

		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.check):
			if m.editing {
				break
			}
			if i, ok := m.list.SelectedItem().(entity.Task); ok {
				i.Done = !i.Done
				m.list.SetItem(m.list.Index(), i)
			}

		case key.Matches(msg, m.keymap.swapUp):
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

			m.list.SetItem(m.list.Index(), cur)
			m.list.SetItem(m.list.Index()+1, above)

		case key.Matches(msg, m.keymap.swapDown):
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

			m.list.SetItem(m.list.Index(), cur)
			m.list.SetItem(m.list.Index()-1, below)

		case key.Matches(msg, m.keymap.insert):
			if m.editing {
				break
			}

			m.list.InsertItem(m.list.Index(), entity.Task{})
			m.list.Select(m.list.Index())

		case key.Matches(msg, m.keymap.remove):
			if m.editing {
				break
			}
			if m.list.Index() == len(m.list.Items())-1 && len(m.list.Items()) > 1 {
				m.list.RemoveItem(m.list.Index())
				m.list.Select(m.list.Index() - 1)
			} else {
				m.list.RemoveItem(m.list.Index())
			}

		case key.Matches(msg, m.keymap.edit):
			if m.editing {
				break
			}
			m.editing = true
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}
			m.textInput.SetValue(cur.Description)

		case key.Matches(msg, m.keymap.editExit):
			if !m.editing {
				break
			}
			cur, ok := m.list.SelectedItem().(entity.Task)
			if !ok {
				break
			}

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
			"\n%s\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}
	return m.list.View()
}

type taskDelegate struct{}

func (d taskDelegate) Height() int                               { return 1 }
func (d taskDelegate) Spacing() int                              { return 0 }
func (d taskDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d taskDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(entity.Task)
	if !ok {
		return
	}

	state := entity.CheckToDo
	if i.Done {
		state = entity.CheckDone
	}
	str := fmt.Sprintf("%s %s", state, i.Description)

	// remove multiple lines
	if idx := strings.Index(str, "\n"); idx != -1 {
		str = str[:idx] + "..."
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("▸ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2)
)

func Run(cfg *config.Config, tasks []entity.Task) []entity.Task {

	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}

	// build the input
	ta := textarea.New()
	ta.Placeholder = "Something todo..."
	ta.Focus()
	ta.CharLimit = 156
	ta.MaxHeight = 30

	l := list.New(
		items,
		taskDelegate{},
		defaultWidth,
		10,
	)

	keyMap := keymap{
		check: key.NewBinding(
			key.WithKeys(cfg.Check),
			key.WithHelp(cfg.Check, "(un)check the tasks"),
		),
		quit: key.NewBinding(
			key.WithKeys(cfg.Quit),
			key.WithHelp(cfg.Quit, "quit"),
		),
		swapUp: key.NewBinding(
			key.WithKeys(cfg.SwapUp),
			key.WithHelp(cfg.SwapUp, "swap up"),
		),
		swapDown: key.NewBinding(
			key.WithKeys(cfg.SwapDown),
			key.WithHelp(cfg.SwapDown, "swap down"),
		),
		remove: key.NewBinding(
			key.WithKeys(cfg.Remove),
			key.WithHelp(cfg.Remove, "remove"),
		),
		insert: key.NewBinding(
			key.WithKeys(cfg.Insert),
			key.WithHelp(cfg.Insert, "insert a new task"),
		),
		edit: key.NewBinding(
			key.WithKeys(cfg.Edit),
			key.WithHelp(cfg.Edit, "edit"),
		),
		editExit: key.NewBinding(
			key.WithKeys(cfg.EditExit),
			key.WithHelp(cfg.EditExit, "edit exit"),
		),
		up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("up", ""),
		),
		down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("down", ""),
		),
	}

	// build the list
	l.Title = ""
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().Height(0).Margin(0, 0, 0, 0).Padding(0, 0, 0, 0)
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(5)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(2)

	m := model{
		list:      l,
		textInput: ta,
		keymap:    keyMap,
	}

	pg := tea.NewProgram(m)
	endmodel, err := pg.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var res []entity.Task
	for _, item := range endmodel.(model).list.Items() {
		res = append(res, item.(entity.Task))
	}
	return res
}
