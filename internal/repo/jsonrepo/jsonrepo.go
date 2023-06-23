package jsonrepo

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/Coding-Brownies/todo/internal/app"
	"github.com/Coding-Brownies/todo/internal/entity"
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

	return res, nil
}

func (j *JSONRepo) Add(t *entity.Task) error {
	tasks, err := j.List()
	if err != nil {
		return err
	}

	tasks = append(tasks, *t)

	err = j.Store(tasks)
	return err
}

func (j *JSONRepo) Store(tasks []entity.Task) error {
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

// funzione locale cucita secondo il contesto e le necessità delle funzioni Delete, Check, Uncheck e Edit
// controlla che un id fornito esista e che la funzione List vada a buon fine, dopodichè restituise index, tasks e un errore
func (j *JSONRepo) idAndListCheck(id string) (int, []entity.Task, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	index := int(i)
	if err != nil {
		return -1, nil, err
	}

	tasks, err := j.List()
	if err != nil {
		return -1, nil, err
	}

	if index < 0 || int(index) >= len(tasks) {
		return -1, nil, errors.New("wrong id")
	}

	return index, tasks, err
}

func (j *JSONRepo) Delete(id string) error {
	index, tasks, err := j.idAndListCheck(id)
	if err != nil {
		return err
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	err = j.Store(tasks)
	return err
}

func (j *JSONRepo) Check(id string) error {
	index, tasks, err := j.idAndListCheck(id)
	if err != nil {
		return err
	}

	tasks[index].Done = true

	err = j.Store(tasks)
	return err
}

func (j *JSONRepo) Uncheck(id string) error {
	index, tasks, err := j.idAndListCheck(id)
	if err != nil {
		return err
	}

	tasks[index].Done = false

	err = j.Store(tasks)
	return err
}

func (j *JSONRepo) Edit(id string, newDescription string) error {
	index, tasks, err := j.idAndListCheck(id)
	if err != nil {
		return err
	}

	tasks[index].Description = newDescription

	err = j.Store(tasks)
	return err
}

func (j *JSONRepo) Swap(IDa string, IDb string) error {
	indexA, tasksA, errA := j.idAndListCheck(IDa)
	if errA != nil {
		return errA
	}
	indexB, tasksB, errB := j.idAndListCheck(IDb)
	if errB != nil {
		return errB
	}
	tasksA[indexA].Position = tasksB[indexB].Position
	tasksB[indexB].Position = tasksA[indexA].Position
	return nil
}
