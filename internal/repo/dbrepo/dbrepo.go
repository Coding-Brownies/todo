package dbrepo

import (
	"encoding/json"
	"errors"

	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	CREATE byte = iota
	DELETE
	UPDATE
)

var _ internal.Repo = &DBRepo{}

type DBRepo struct {
	DB *gorm.DB
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
	err := db.DB.Order("position").Find(&res).Error
	return res, err
}

func (db *DBRepo) Add(t *entity.Task) error {
	t.ID = uuid.NewString()
	err := db.DB.Create(t).Error
	if err != nil {
		return err
	}
	// Crea e registra l'azione add nella tabella di registro delle modifiche Change
	return db.do(t, CREATE, uuid.NewString())
}

func (db *DBRepo) Delete(t *entity.Task) error {
	// Crea e registra l'azione delete nella tabella di registro delle modifiche Change
	err := db.do(t, DELETE, uuid.NewString())
	if err != nil {
		return err
	}
	// elimina il task
	return db.DB.Where("id=?", t.ID).Delete(&entity.Task{}).Error
}

func (db *DBRepo) Check(t *entity.Task) error {
	// Crea e registra l'azione check nella tabella di registro delle modifiche Change
	err := db.do(t, UPDATE, uuid.NewString())
	if err != nil {
		return err
	}
	return db.DB.Model(&entity.Task{}).Where("id = ?", t.ID).Update("done", true).Error
}

func (db *DBRepo) Uncheck(t *entity.Task) error {
	// Crea e registra l'azione uncheck nella tabella di registro delle modifiche Change
	err := db.do(t, UPDATE, uuid.NewString())
	if err != nil {
		return err
	}
	return db.DB.Model(&entity.Task{}).Where("id = ?", t.ID).Update("done", false).Error
}

func (db *DBRepo) Edit(t *entity.Task, newDescription string) error {
	// Crea e registra l'azione edit nella tabella di registro delle modifiche Change
	err := db.do(t, UPDATE, uuid.NewString())
	if err != nil {
		return err
	}
	// edit del campo description
	return db.DB.Model(&entity.Task{}).Where("id = ?", t.ID).Update("description", newDescription).Error
}

func (db *DBRepo) Swap(taskA, taskB *entity.Task) error {
	actionID := uuid.New().String()
	err := db.do(taskA, UPDATE, actionID)
	if err != nil {
		return err
	}
	err = db.do(taskB, UPDATE, actionID)
	if err != nil {
		return err
	}
	// effettuo la swap
	err = db.DB.Model(&entity.Task{}).Where("id=?", taskA.ID).Update("position", taskB.Position).Error
	if err != nil {
		return err
	}
	err = db.DB.Model(&entity.Task{}).Where("id=?", taskB.ID).Update("position", taskA.Position).Error
	if err != nil {
		return err
	}
	return nil
}

// funzione ausiliaria Do che accetta un task,
// ne esegue il marshalling e lo salva nelle tabella di registro delle modifiche Change
func (db *DBRepo) do(task *entity.Task, action byte, actionID string) error {
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
	return db.DB.Create(&change).Error
}

func (db *DBRepo) Undo() error {
	var changes []entity.Change
	// find the action_id of the last change
	subQuery := db.DB.Select("action_id").Order("id desc").Table("changes").Limit(1)
	// recupera tutte le changes aventi quell'actionID
	err := db.DB.Where("action_id = (?)", subQuery).Find(&changes).Error
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		return errors.New("last change not found")
	}
	// delete the changes aventi quell'action_id
	err = db.DB.Where("action_id = ?", changes[0].ActionID).Delete(&entity.Change{}).Error
	if err != nil {
		return err
	}

	for _, change := range changes {
		// per accedere alle informazioni del task nel campo OldStatus, decodificare i byte JSON utilizzando json.Unmarshal()
		var oldStatus entity.Task
		err = json.Unmarshal(change.OldStatus, &oldStatus)
		if err != nil {
			return err
		}
		// Effettua il revert dell'azione
		switch change.Action {
		case CREATE:
			err = db.DB.Where("id=?", oldStatus.ID).Delete(&entity.Task{}).Error
		case DELETE, UPDATE:
			err = db.DB.Save(&oldStatus).Error
		default:
			err = errors.New("unsupported action")
		}
		if err != nil {
			return err
		}
	}
	return nil
}
