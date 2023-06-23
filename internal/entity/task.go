// create a folder entity with a file for each used entity (ex: task.go which is a struct)
package entity

import "time"

const CheckToDo = "◻"
const CheckDone = "◼"

type Task struct {
	ID          string
	Done        bool
	Description string
	Position    time.Time `gorm:"autoCreateTime"`
}

func (t Task) FilterValue() string { return "" }
