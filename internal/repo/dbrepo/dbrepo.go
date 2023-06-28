package dbrepo

import (
	"encoding/json"
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

// struct per la tabella registro delle modifiche
type ChangeLog struct {
	ID        string      `json:"id"`
	Action    string      `json:"action"` // Nome della funzione interessata
	Data      interface{} `json:"data"`
	CreatedAt time.Time   // Data e ora di creazione del log
}

// struttura dati che rappresenta l'oggetto JSON da memorizzare nel campo "data" di ChangeLog
// AddData contiene le informazioni necessarie da salvare per ogni tipo di azione "Add"
// i tag json servono per specificare i nomi dei campi come verranno serializzati e deserializzati nel formato JSON
type CheckData struct {
	ID string `json:"id"`
}

// Funzioni ausiliare della struct ChangeLog:

// Funzione per registrare una modifica c nella tabella delle modifiche
func (db *DBRepo) RegisterChange(c ChangeLog) error {
	err := db.Create(&c).Error
	return err
}

// funzione per ottenere una modifica dalla tabella di registro delle modifiche
// il metodo First seleziona la prima riga dalla tabella che corrisponde all'ID fornito
// il risultato viene memorizzato nella variabile change di tipo ChangeLog
// Se l'operazione ha successo, la funzione restituisce la modifica trovata
// In caso di errore, viene restituito un oggetto vuoto di tipo ChangeLog e l'errore associato
func (db *DBRepo) GetChange(changeID string) (ChangeLog, error) {
	var change ChangeLog
	err := db.Table("RegisterChange").Where("id = ?", changeID).
		First(&change).Error
	if err != nil {
		return ChangeLog{}, err
	}

	var data interface{}
	switch change.Action {
	case "Check":
		var checkData CheckData
		err = json.Unmarshal(change.Data, &checkData)
		data = checkData

		// il resto dei casi per le altre azioni
	}
	// Restituisci la modifica con i dati decodificati
	return ChangeLog{
		ID:     change.ID,
		Action: change.Action,
		Data:   data,
	}, nil

}

// funzione che elimina una modifica dalla tabella di registro delle modifiche
// il metodo Delete elimina dalla tabella la riga che corrisponde all'ID fornito
// Il parametro &ChangeLog{} viene passato come destinazione per indicare quale struttura di dati deve essere eliminata
func (db *DBRepo) DeleteChange(changeID string) error {
	return db.Table("RegisterChange").Where("id = ?", changeID).
		Delete(&ChangeLog{}).Error
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
	err := db.Model(&entity.Task{}).Where("id=?", ID).
		Update("done", true).Error
	if err != nil {
		return err
	}
	// registrazione della modifica appena effettuata
	data := CheckData{
		ID: ID,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	change := ChangeLog{
		Action: "Check",
		Data:   jsonData,
	}

	return db.RegisterChange(change)
}

func (db *DBRepo) Uncheck(ID string) error {
	return db.Model(&entity.Task{}).Where("id=?", ID).
		Update("done", false).Error
}

func (db *DBRepo) Edit(ID string, newDescription string) error {
	return db.Model(&entity.Task{}).Where("id=?", ID).
		Update("description", newDescription).Error
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

// Funzione che esegue il rollback di una modifica utilizzando la tabella di registro delle modifihe
// Undo prende in ingresso l'ID della modifica da annullare, e deve cercare la modifica nella tabella di registro delle modifiche
// e chiamare la funzione opposta per annullare l'azione
// Add => Delete e viceversa, Edit => Edit
func (db *DBRepo) Undo(changeID string) error {
	change, err := db.GetChange(changeID)
	if err != nil {
		return err
	}
	// Chiamare la funzione opposta in base all'azione registrata
	switch change.Action {
	case "Add":
		err = db.Delete(change.Data.(*entity.Task).ID)
	case "Delete":
		err = db.Add(change.Data.(*entity.Task))
	case "Check":
		var checkData CheckData
		err = json.Unmarshal(change.Data, &checkData)
		if err != nil {
			return err
		}
		// Chiamare la funzione Uncheck utilizzando l'ID decodificato
		err = db.Uncheck(checkData.ID)
	case "Uncheck":
		err = db.Check(change.Data.(string))
	case "Edit":
		data := change.Data.(map[string]interface{})
		err = db.Edit(data["ID"].(string), data["oldDescription"].(string))
	case "Swap":
		data := change.Data.(map[string]interface{})
		err = db.Swap(data["IDa"].(string), data["IDb"].(string))
	}
	// gestione degli errori dei case
	if err != nil {
		return err
	}
	// Elimina la modifica dalla tabella di registro delle modifiche
	return db.DeleteChange(changeID)
}
