package dbrepo

import (
	"fmt"

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

func New(dbpath string) (*DBRepo, error) {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&entity.Task{}, &entity.Change{})

	return &DBRepo{
		DB: db,
	}, nil
}

func (db *DBRepo) List() ([]entity.Task, error) {
	var res []entity.Task
	// print while ordering by Position (type time.Time)
	err := db.Order("position").Find(&res).Error
	return res, err
}

func (db *DBRepo) Add(t *entity.Task) error {
	t.ID = uuid.New().String()
	err := db.Create(t).Error
	if err != nil {
		return err
	}
	// Registra l'azione di aggiunta nella tabella di registro delle modifiche
	change := entity.Change{
		TaskID:      t.ID,
		Action:      "Add",
		Description: t.Description,
		Position:    t.Position,
	}
	return db.Create(&change).Error
}

func (db *DBRepo) Delete(ID string) error {
	// Recupera il task prima di eliminarlo
	var task entity.Task
	err := db.First(&task, "id = ?", ID).Error
	if err != nil {
		return err
	}
	// elimina il task
	err = db.Where("id=?", ID).Delete(&entity.Task{}).Error
	if err != nil {
		return err
	}
	// Registra l'azione di eliminazione nella tabella di registro delle modifiche
	change := entity.Change{
		TaskID:      ID,
		Action:      "Delete",
		Description: task.Description,
		Position:    task.Position,
	}
	return db.Create(&change).Error

}

func (db *DBRepo) Check(ID string) error {
	err := db.Model(&entity.Task{}).Where("id = ?", ID).Update("done", true).Error
	if err != nil {
		return err
	}
	// Registra l'azione di check nella tabella di registro delle modifiche
	change := entity.Change{
		TaskID: ID,
		Action: "Check",
	}
	return db.Create(&change).Error
}

func (db *DBRepo) Uncheck(ID string) error {
	err := db.Model(&entity.Task{}).Where("id = ?", ID).Update("done", false).Error
	if err != nil {
		return err
	}
	// Registra l'azione di uncheck nella tabella di registro delle modifiche
	change := entity.Change{
		TaskID: ID,
		Action: "Uncheck",
	}
	return db.Create(&change).Error
}

func (db *DBRepo) Edit(ID string, newDescription string) error {
	old := &entity.Task{}
	err := db.Where("id = ?", ID).First(old).Error
	if err != nil {
		return err
	}
	change := entity.Change{
		TaskID:      ID,
		Action:      "Edit",
		Description: old.Description,
	}

	err = db.Model(&entity.Task{}).Where("id = ?", ID).Update("description", newDescription).Error
	if err != nil {
		return err
	}
	// Registra l'azione di modifica nella tabella di registro delle modifiche

	return db.Create(&change).Error
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
	// Registra l'azione di swap nella tabella di registro delle modifiche
}

func (db *DBRepo) Undo() error {
	var change entity.Change
	// Trova l'ultima modifica registrata non ancora revertita
	err := db.Where("reverted = ?", false).Order("id desc").First(&change).Error
	if err != nil {
		return err
	}
	// Aggiorna il flag di revert nella tabella di registro delle modifiche
	err = db.Model(&entity.Change{}).Where("id = ?", change.ID).Update("reverted", true).Error
	if err != nil {
		return err
	}
	// Effettua il revert dell'azione
	switch change.Action {
	case "Add":
		return db.Delete(change.TaskID)
	case "Delete":
		// Ripristina il task eliminato
		task := entity.Task{
			ID:          change.TaskID,
			Description: change.Description,
			Position:    change.Position,
		}
		return db.Create(&task).Error
	case "Check":
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("done", false).Error
	case "Uncheck":
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("done", true).Error
	case "Edit":
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("description", change.Description).Error
	//case "Swap":
	//return db.Swap()
	default:
		return fmt.Errorf("azione non supportata: %s", change.Action)
	}
}
