package bubble

import (
	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble/components/bin"
	"github.com/Coding-Brownies/todo/internal/bubble/components/edit"
	"github.com/Coding-Brownies/todo/internal/bubble/components/task"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &model{}

type ModelWithHelp interface {
	tea.Model
	help.KeyMap
}

type model struct {
	repo   internal.Repo
	helper help.Model
	// components
	tasks *task.Model
	bin   *bin.Model

	// states
	cur     ModelWithHelp
	bigHelp bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m.cur.Update(msg)
}

func (m model) View() string {
	view := "\n" + m.cur.View()
	help := m.helper.ShortHelpView(m.cur.ShortHelp())

	if m.bigHelp {
		help = m.helper.FullHelpView(m.cur.FullHelp())
	}

	return view +
		list.DefaultStyles().
			HelpStyle.PaddingLeft(2).
			Foreground(lipgloss.Color("#000000")).Render(help)
}

func New(cfg *config.Config, repo internal.Repo) *model {

	keyMap := NewKeyMap(cfg)

	return &model{
		repo: repo,
		tasks: task.NewModel(
			task.KeyMap{
				Check:    keyMap.Check,
				Quit:     keyMap.Quit,
				SwapUp:   keyMap.SwapUp,
				SwapDown: keyMap.SwapDown,
				Remove:   keyMap.Remove,
				Insert:   keyMap.Insert,
				Up:       keyMap.Up,
				Down:     keyMap.Down,
			},
			repo,
			edit.NewModel(
				edit.KeyMap{
					Exit: keyMap.EditExit,
				},
				repo,
			),
		),
		bin: bin.NewModel(
			bin.KeyMap{
				Up:       keyMap.Up,
				Down:     keyMap.Down,
				Restore:  keyMap.Restore,
				EmptyBin: keyMap.EmptyBin,
			},
			repo,
		),
	}
}

func (m *model) Run() error {
	tasks, err := m.repo.List()
	if err != nil {
		return err
	}
	m.tasks.Fill(tasks...)
	m.cur = m.tasks

	pg := tea.NewProgram(m)
	_, err = pg.Run()
	if err != nil {
		return err
	}
	return nil
}
