package dbrepo

import (
	"time"

	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ internal.Repo = &DBRepo{}

type DBRepo struct {
	*gorm.DB
}

// struct per la tabella delle modifiche
type ChangeLog struct {
	ID        uint      `gorm:"primaryKey"`
	FuncName  string    // Nome della funzione interessata
	TaskID    string    // ID del task campo interessato
	FieldName string    // Nome del campo interessato
	OldValue  string    // Valore originale prima della modifica
	CreatedAt time.Time // Data e ora di creazione del log
}

func New(dbpath string) (*DBRepo, error) {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&entity.Task{})

	return &DBRepo{
		DB: db,
	}, nil
}

func (db *DBRepo) List() ([]entity.Task, error) {
	var res []entity.Task
	// print while ordering by Position (type time.Time)

	err := db.Order("position").Find(&res).Error

	// in caso di errore res è vuoto
	return res, err
}

func (db *DBRepo) Add(t *entity.Task) error {
	t.ID = uuid.New().String()
	return db.Create(t).Error
}

func (db *DBRepo) Delete(ID string) error {
	return db.Where("id=?", ID).Delete(&entity.Task{}).Error
}

func (db *DBRepo) Check(ID string) error {
	return db.Model(&entity.Task{}).Where("id=?", ID).Update("done", true).Error
}

func (db *DBRepo) Uncheck(ID string) error {
	return db.Model(&entity.Task{}).Where("id=?", ID).Update("done", false).Error
}

func (db *DBRepo) Edit(ID string, newDescription string) error {
	// Registrare il valore originale prima della modifica
	var oldTask entity.Task
	err := db.Model(&entity.Task{}).Where("id = ?", ID).First(&oldTask).Error
	if err != nil {
		return err
	}

	// Registra la modifica nella tabella delle modifiche
	changeLog := ChangeLog{
		FuncName:  "edit",              // Nome della funzione delle modifiche
		TaskID:    ID,                  // ID del task interessato
		FieldName: "description",       // Nome del campo modificato
		OldValue:  oldTask.Description, // Valore originale prima della modifica
		CreatedAt: time.Now(),          // Data e ora di creazione del log
	}
	err = db.Create(&changeLog).Error
	if err != nil {
		return err
	}

	// effettua la modifica effettiva
	return db.Model(&entity.Task{}).Where("id=?", ID).Update("description", newDescription).Error
}

func (db *DBRepo) Swap(IDa string, IDb string) error {
	var a, b entity.Task

	err := db.Select("position").Take(&a, "id=?", IDa).Error
	if err != nil {
		return err
	}
	err = db.Select("position").Take(&b, "id=?", IDb).Error
	if err != nil {
		return err
	}

	err = db.Model(&entity.Task{}).Where("id=?", IDa).Update("position", b.Position).Error
	if err != nil {
		return err
	}
	err = db.Model(&entity.Task{}).Where("id=?", IDb).Update("position", a.Position).Error
	if err != nil {
		return err
	}

	return nil
}

// Funzione per eseguire il rollback di una modifica utilizzando la tabella delle modifiche
func (db *DBRepo) UndoLastChange(ID string) error {
	// Recupera l'ultima modifica per il record e il campo specificati
	var changeLog ChangeLog
	err := db.Where("funcname = ? AND taskid = ?", "edit", ID).Order("created_at DESC").First(&changeLog).Error
	if err != nil {
		return err
	}

	// Effettua il rollback ripristinando il valore originale nel campo description
	err = db.Model(&entity.Task{}).Where("id = ?", ID).
		Update("description", changeLog.OldValue).Error
	if err != nil {
		return err
	}

	// Elimina il log di modifica dalla tabella delle modifiche
	return nil // return db.DeleteChangeLog(&changeLog)
}

// funzione delete per il tipo changelog, perè prende la delete sbagliata
/*
func (db *DBRepo) DeleteChangeLog(changeLog *ChangeLog) error {
	return db.Delete(changeLog).Error
}
*/
