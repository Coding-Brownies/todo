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
