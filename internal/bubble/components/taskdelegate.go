package components

import (
	"fmt"
	"io"
	"strings"

	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2)
)

// this struct is responsible for the rendering of an item inside the list
type CustomItemRender struct{}

func (d CustomItemRender) Height() int                               { return 1 }
func (d CustomItemRender) Spacing() int                              { return 0 }
func (d CustomItemRender) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d CustomItemRender) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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
