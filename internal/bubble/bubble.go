package bubble

import (
	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble/components/bin"
	"github.com/Coding-Brownies/todo/internal/bubble/components/edit"
	"github.com/Coding-Brownies/todo/internal/bubble/components/task"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &model{}

type model struct {
	repo internal.Repo
	// components
	tasks *task.Model
	bin   *bin.Model
	edit  *edit.Model

	cur tea.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m model) View() string {

	return m.cur.View()

	// help := m.list.Help.ShortHelpView(m.keymap.ShortHelp())
	// if m.bigHelp {
	// 	help = m.list.Help.FullHelpView(m.keymap.FullHelp())
	// }
	// return "\n" + m.list.View() + m.list.Styles.HelpStyle.Render(help)
}

func New(cfg *config.Config, repo internal.Repo) *model {

	keyMap := NewKeyMap(cfg)

	return &model{
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
