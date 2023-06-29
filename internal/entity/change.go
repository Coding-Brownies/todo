package entity

import "time"

// struct per la tabella registro delle modifiche
type Change struct {
	ID uint `gorm:"primaryKey,autoIncrement"`

	Action      string
	TaskID      string
	Description string
	Position    time.Time

	Reverted bool
	//RevertedByID uint
	//RevertedAt   time.Time
}
