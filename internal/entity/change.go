package entity

import "time"

// struct per la tabella di registro delle modifiche
type Change struct {
	ID uint `gorm:"primaryKey,autoIncrement"`

	Action      string
	TaskID      string
	Description string
	Position    time.Time

	ActionID string // id dell'operazione di modifica, necessario per le operazioni che richiedono 2 id (es: swap)
}
