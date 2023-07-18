package edit

import (
	"fmt"

	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/bubble/components/task"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

var _ bubble.BubbleModel = &Model{}
var _ task.Editor = &Model{}

type Model struct {
	textarea.Model

	keymap *KeyMap
	err    error
	Res    string
}

func (m *Model) Map() help.KeyMap {
	return m.keymap
}

// Error implements bubble.BubbleModel.
func (m *Model) Error() error {
	return m.err
}

// Edit implements task.Editor.
func (m *Model) Edit(string) tea.Cmd {
	return func() tea.Msg {

		pg := tea.NewProgram(m)
		res, err := pg.Run()
		if err != nil {
			return err
		}

		if text, ok := res.(Model); ok {
			return task.EditFinished{
				Content: text.Value(),
			}
		}

		return nil
	}
}

func NewModel(k *KeyMap) *Model {
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
		m.err = msg
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
