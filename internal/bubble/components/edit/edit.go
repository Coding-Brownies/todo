package edit

import (
	"errors"
	"fmt"

	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble/components/task"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Model{}
var _ task.Editor = &Model{}

type Model struct {
	textarea.Model

	keymap KeyMap
	Error  error
	Res    string
}

// Edit implements task.Editor.
func (m *Model) Edit(string) (string, error) {

	pg := tea.NewProgram(m)
	resultModel, err := pg.Run()
	if err != nil {
		return "", err
	}

	if model, ok := resultModel.(Model); ok {
		return model.Value(), nil
	}

	// TODO: make this better
	return "", errors.New("boh")
}

func NewModel(k KeyMap, r internal.Repo) *Model {
	ta := textarea.New()
	ta.Placeholder = "Something todo..."
	ta.Focus()
	ta.CharLimit = 156
	ta.MaxHeight = 30

	return &Model{
		Model:  ta,
		keymap: k,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keymap.Exit) {
			return m, tea.Quit
		}
	case error:
		m.Error = msg
		return m, nil
	}

	_, cmd := m.Model.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m Model) View() string {
	return fmt.Sprintf(
		"\n  🖊️  Edit\n\n%s\n\n",
		m.Model.View(),
	)
}
