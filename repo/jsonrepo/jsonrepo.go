package jsonrepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Coding-Brownies/todo/app"
	"github.com/Coding-Brownies/todo/entity"
)

var _ app.Repo = &JSONRepo{}

type JSONRepo struct {
	path string
}

// the New function of the  JSONRepo will accepts a path to the .json file
func New(p string) *JSONRepo {
	return &JSONRepo{
		path: p,
	}
}

func (j *JSONRepo) List() ([]entity.Task, error) {
	// creare il file json se non esiste già

	// check if file exists
	_, err := os.Stat(j.path)
	// if file does not exist, create file
	if err != nil {
		_, err := os.Create(j.path)

		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(j.path)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var res []entity.Task

	if string(content) != "" {
		if err := json.Unmarshal(content, &res); err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(res); i++ {
		res[i].ID = fmt.Sprint(i)
	}
	return res, nil
}

func (j *JSONRepo) Add(t *entity.Task) error {
	tasks, err := j.List()
	if err != nil {
		return err
	}
	tasks = append(tasks, *t)
	err = j.store(tasks)
	return err
}

// funzione in locale
func (j *JSONRepo) store(tasks []entity.Task) error {
	// conversione di tasks in JSON
	content, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	// scrittura del JSON nel file, sovrascrivendo l'eventuale contenuto precedente se il file non è vuoto
	// il codice 0644 da il permesso per sovrascirvere il file
	err = os.WriteFile(j.path, content, 0644)
	// in caso di errore durante la conversione o la scrittura, la funzione restituisce un errore
	// nel caso in cui non ci siano stati errori err è nill
	return err
}

func (j *JSONRepo) Delete(id string) error {
	index, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	tasks, err := j.List()
	if err != nil {
		return err
	}
	if index < 0 || int(index) >= len(tasks) {
		return errors.New("wrong id")
	}
	tasks = append(tasks[:index], tasks[index+1:]...)
	err = j.store(tasks)
	return err
}

func (j *JSONRepo) Check(id string) error {
	index, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	tasks, err := j.List()
	if err != nil {
		return err
	}
	if index < 0 || int(index) >= len(tasks) {
		return errors.New("wrong id")
	}

	tasks[index].Done = true
	err = j.store(tasks)
	return err
}

func (j *JSONRepo) Uncheck(id string) error {
	index, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	tasks, err := j.List()
	if err != nil {
		return err
	}
	if index < 0 || int(index) >= len(tasks) {
		return errors.New("wrong id")
	}

	tasks[index].Done = false
	err = j.store(tasks)
	return err
}
