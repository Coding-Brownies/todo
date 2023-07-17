package components

import (
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/list"
)

func SetList(tasks []entity.Task, l *list.Model) {
	items := make([]list.Item, len(tasks))
	for i, v := range tasks {
		items[i] = v
	}
	l.SetItems(items)
	if l.Index() > len(items)-1 {
		l.Select(len(items) - 1)
	}
	if l.Index() == -1 {
		l.Select(0)
	}
}
