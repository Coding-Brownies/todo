package dbrepo

import (
	"fmt"

	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ internal.Repo = &DBRepo{}

type DBRepo struct {
	*gorm.DB
}

func New(dbpath string) (*DBRepo, error) {
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{
		//TODO: silent only record not found
		Logger: logger.Default.LogMode(logger.Silent),
	})
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
	// Una volta che la modifica è andata a buon fine, registra l'azione di modifica nella tabella di registro delle modifiche
	return db.Create(&change).Error
}

func (db *DBRepo) Swap(taskA, taskB *entity.Task) error {
	// creazione delle due change, con il valore attuale del campo position (prima di effettuare la swap)
	actionID := uuid.New().String()
	change := []entity.Change{
		{
			TaskID:   taskA.ID,
			Action:   "Swap",
			ActionID: actionID,
			Position: taskA.Position,
		},
		{
			TaskID:   taskB.ID,
			Action:   "Swap",
			ActionID: actionID,
			Position: taskB.Position,
		},
	}
	// effettuo la swap
	err := db.Model(&entity.Task{}).Where("id=?", taskA.ID).Update("position", taskB.Position).Error
	if err != nil {
		return err
	}
	err = db.Model(&entity.Task{}).Where("id=?", taskB.ID).Update("position", taskA.Position).Error
	if err != nil {
		return err
	}
	// Registra le due change nella tabella di registro delle modifiche
	err = db.Create(&change[0]).Error
	if err != nil {
		return err
	}
	err = db.Create(&change[1]).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DBRepo) Undo() error {
	var change entity.Change
	// Trova l'ultima modifica registrata
	err := db.Order("id desc").First(&change).Error
	if err != nil {
		return err
	}
	// e la elimina la change
	err = db.Where("id = ?", change.ID).Delete(&entity.Change{}).Error
	if err != nil {
		return err
	}
	// Effettua il revert dell'azione
	switch change.Action {
	case "Add":
		return db.DB.Where("id=?", change.TaskID).Delete(&entity.Task{}).Error
	case "Delete":
		// Ripristina il task eliminato
		task := entity.Task{
			ID:          change.TaskID,
			Description: change.Description,
			Position:    change.Position,
		}
		return db.DB.Create(&task).Error
	case "Check":
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("done", false).Error
	case "Uncheck":
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("done", true).Error
	case "Edit":
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("description", change.Description).Error
	case "Swap":
		// prelevare il secondo change legato all'action id del primo change
		var changeB entity.Change
		err := db.Where("action_id=?", change.ActionID).First(&changeB).Error //first , ma a regola funziona anche con find
		if err != nil {
			return err
		}
		// swap in query
		err = db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("position", change.Position).Error
		if err != nil {
			return err
		}
		err = db.Model(&entity.Task{}).Where("id = ?", changeB.TaskID).Update("position", changeB.Position).Error
		if err != nil {
			return err
		}
		// cancello il secondo change
		err = db.Where("id = ?", changeB.ID).Delete(&entity.Change{}).Error
		if err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("azione non supportata: %s", change.Action)
	}
}
