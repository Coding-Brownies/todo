package bin

import (
	"fmt"
	"io"
	"strings"

	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Faint(true)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Faint(true)
)

// this struct is responsible for the rendering of an item inside the list
type CustomItemRender struct {
	Editing *bool
	Editor  *textinput.Model
}

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

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("▸ " + strings.Join(s, " "))
		}
	}

	if d.Editing != nil && *d.Editing && index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("  "+state+" "+d.Editor.View()))
		return
	}

	fmt.Fprint(w, fn(str))
}
