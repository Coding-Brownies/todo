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
	actionID := uuid.New().String()
	return db.do(t, "Create", actionID)
}

func (db *DBRepo) Delete(t *entity.Task) error {
	// Crea e registra l'azione delete nella tabella di registro delle modifiche Change
	actionID := uuid.New().String()
	err := db.do(t, "Delete", actionID)
	if err != nil {
		return err
	}
	// elimina il task
	return db.Where("id=?", t.ID).Delete(&entity.Task{}).Error
}

func (db *DBRepo) Check(t *entity.Task) error {
	// Crea e registra l'azione check nella tabella di registro delle modifiche Change
	actionID := uuid.New().String()
	err := db.do(t, "Update", actionID)
	if err != nil {
		return err
	}
	return db.Model(&entity.Task{}).Where("id = ?", t.ID).Update("done", true).Error
}

func (db *DBRepo) Uncheck(t *entity.Task) error {
	// Crea e registra l'azione uncheck nella tabella di registro delle modifiche Change
	actionID := uuid.New().String()
	err := db.do(t, "Update", actionID)
	if err != nil {
		return err
	}
	return db.Model(&entity.Task{}).Where("id = ?", t.ID).Update("done", false).Error
}

func (db *DBRepo) Edit(t *entity.Task, newDescription string) error {
	// Crea e registra l'azione edit nella tabella di registro delle modifiche Change
	actionID := uuid.New().String()
	err := db.do(t, "Update", actionID)
	if err != nil {
		return err
	}
	// edit del campo description
	return db.Model(&entity.Task{}).Where("id = ?", t.ID).Update("description", newDescription).Error
}

func (db *DBRepo) Swap(taskA, taskB *entity.Task) error {
	actionID := uuid.New().String()
	err := db.do(taskA, "Update", actionID)
	if err != nil {
		return err
	}
	err = db.do(taskB, "Update", actionID)
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
func (db *DBRepo) do(task *entity.Task, action, actionID string) error {
	// codifica del json come []byte usando Marshal
	oldStatusJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
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
	// Trova l'ultima change registrata
	var c entity.Change
	err := db.DB.Order("id desc").First(&c).Error
	if err != nil {
		return err
	}
	// Recupera le changes aventi quell'actionID
	var changes []entity.Change
	err = db.DB.Where("action_id=?", c.ActionID).Find(&changes).Error
	if err != nil {
		return err
	}
	// per accedere alle informazioni del task nel campo OldStatus, decodificare i byte JSON utilizzando json.Unmarshal()
	var oldStatus entity.Task
	err = json.Unmarshal(changes[0].OldStatus, &oldStatus)
	if err != nil {
		return err
	}
	// elimina le changes correlate a quell'actionID, dato che ne viene effettuato il revert
	err = db.DB.Exec("DELETE FROM changes WHERE action_id = ?", changes[0].ActionID).Error
	if err != nil {
		return err
	}
	// Effettua il revert dell'azione
	switch changes[0].Action {
	case "Create":
		return db.DB.Where("id=?", oldStatus.ID).Delete(&entity.Task{}).Error
	case "Delete":
		return db.DB.Save(&oldStatus).Error
	case "Update":
		// avrà lunghezza 2 solo nel caso dello swap, trattare come 2 Update
		if len(changes) == 2 {
			var oldStatusB entity.Task
			err = json.Unmarshal(changes[1].OldStatus, &oldStatusB)
			if err != nil {
				return err
			}
			err = db.DB.Save(&oldStatusB).Error
			if err != nil {
				return err
			}
		}
		// casi check, uncheck, edit
		return db.DB.Save(&oldStatus).Error
	default:
		return fmt.Errorf("azione non supportata: %s", changes[0].Action)
	}
}
