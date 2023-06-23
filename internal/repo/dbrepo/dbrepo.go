package dbrepo

import (
	"github.com/Coding-Brownies/todo/internal/app"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ app.Repo = &DBRepo{}

type DBRepo struct {
	*gorm.DB
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
	err := db.Find(&res).Error
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

func (db *DBRepo) Store(tasks []entity.Task) error {
	for i := 0; i < len(tasks); i++ {
		tasks[i].ID = uuid.New().String()
	}
	err := db.Where("1=1").Delete(&entity.Task{}).Error
	if err != nil {
		return err
	}
	return db.Create(tasks).Error
}

func (db *DBRepo) Edit(ID string, newDescription string) error {
	return db.Model(&entity.Task{}).Where("id=?", ID).Update("description", newDescription).Error
}
