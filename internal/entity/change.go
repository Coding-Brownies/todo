package entity

import "time"

// struct per la tabella registro delle modifiche
type Change struct {
	ID uint `gorm:"primaryKey,autoIncrement"`

	Action      string
	TaskID      string
	Description string
	Position    time.Time

	ActionID string // id dell'operazione, necessario per le operazioni che richiedono 2 id (swap)
}
