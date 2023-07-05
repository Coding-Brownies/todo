package entity

// struct per la tabella di registro delle modifiche
type Change struct {
	ID uint `gorm:"primaryKey,autoIncrement"`

	Action    byte // usare le keywords Update, Delete o Create
	OldStatus []byte // informazioni del task modificato (ID, Description, Done, Position)

	ActionID string // id dell'operazione di modifica, necessario per le operazioni che richiedono 2 id (es: swap)
}
