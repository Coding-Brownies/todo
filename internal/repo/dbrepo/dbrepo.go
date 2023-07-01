package dbrepo

import (
	"encoding/json"
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
	// Crea e registra l'azione add nella tabella di registro delle modifiche Change
	return db.Do(t, "Create")
}

func (db *DBRepo) Delete(t *entity.Task) error {
	// Crea e registra l'azione delete nella tabella di registro delle modifiche Change
	err := db.Do(t, "Delete")
	if err != nil {
		return err
	}
	// elimina il task
	return db.Where("id=?", t.ID).Delete(&entity.Task{}).Error
}

func (db *DBRepo) Check(t *entity.Task) error {
	// Crea e registra l'azione check nella tabella di registro delle modifiche Change
	err := db.Do(t, "Update")
	if err != nil {
		return err
	}
	return db.Model(&entity.Task{}).Where("id = ?", t.ID).Update("done", true).Error
}

func (db *DBRepo) Uncheck(t *entity.Task) error {
	// Crea e registra l'azione uncheck nella tabella di registro delle modifiche Change
	err := db.Do(t, "Update")
	if err != nil {
		return err
	}
	return db.Model(&entity.Task{}).Where("id = ?", t.ID).Update("done", false).Error
}

func (db *DBRepo) Edit(t *entity.Task, newDescription string) error {
	// Crea e registra l'azione edit nella tabella di registro delle modifiche Change
	err := db.Do(t, "Update")
	if err != nil {
		return err
	}
	// edit del campo description
	return db.Model(&entity.Task{}).Where("id = ?", t.ID).Update("description", newDescription).Error
}

func (db *DBRepo) Swap(taskA, taskB *entity.Task) error {
	err := db.Do(taskA, "Update")
	if err != nil {
		return err
	}
	err = db.Do(taskB, "Update")
	if err != nil {
		return err
	}
	// effettuo la swap
	err = db.Model(&entity.Task{}).Where("id=?", taskA.ID).Update("position", taskB.Position).Error
	if err != nil {
		return err
	}
	err = db.Model(&entity.Task{}).Where("id=?", taskB.ID).Update("position", taskA.Position).Error
	if err != nil {
		return err
	}
	return nil
}

// funzione ausiliaria Do che accetta un task,
// ne esegue il marshalling e lo salva nelle tabella di registro delle modifiche Change
func (db *DBRepo) Do(task *entity.Task, action string) error {
	// codifica del json come []byte usando Marshal
	oldStatusJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
	actionID := uuid.New().String()
	// crea una nuova change
	change := entity.Change{
		Action:    action,
		OldStatus: oldStatusJSON,
		ActionID:  actionID,
	}
	// salva il change nella tabella di registro delle modifiche
	return db.Create(&change).Error
}

func (db *DBRepo) Undo() error {
	// prelevare l'ultimo action id e fare la query per recuperare l'elenco di change con quell'actionID
	var c entity.Change
	// Trova l'ultima modifica registrata
	err := db.Order("id desc").First(&c).Error
	if err != nil {
		return err
	}
	var change []entity.Change
	err = db.Where("action_id=?", c.ActionID).Find(&change).Error
	if err != nil {
		return err
	}
	if len(change) != 1 || len(change) != 2 {
		fmt.Errorf("numero di change collegati all'azione non supportato: %s", len(change))
	}
	// elimina le change correlate a quell'azione
	err = db.Where("action_id = ?", c.ActionID).Delete(&entity.Change{}).Error
	if err != nil {
		return err
	}
	// Effettua il revert dell'azione
	switch change[0].Action {
	case "Create":
		return db.DB.Where("id=?", change[0].TaskID).Delete(&entity.Task{}).Error
	case "Delete":
		// Ripristina il task eliminato
		task := entity.Task{
			ID:          change[0].TaskID,
			Description: change[0].Description,
			Position:    change[0].Position,
		}
		return db.DB.Create(&task).Error
	case "Update":
		// avrà lunghezza 2 solo nel caso dello swap, trattare come 2 Update
		if len(change) == 2 {

		}
		// altri casi di update (check, uncheck, edit)
		return db.Model(&entity.Task{}).Where("id = ?", change.TaskID).Update("done", false).Error
	/*
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
	*/
	default:
		return fmt.Errorf("azione non supportata: %s", change[0].Action)
	}
}
