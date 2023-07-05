package entity

import (
	"time"

	"gorm.io/gorm"
)

const CheckToDo = "◻"
const CheckDone = "◼"

type Task struct {
	ID          string
	Done        bool
	Description string
	Position    time.Time      `gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (t Task) FilterValue() string { return "" }
