package dbrepo_test

import (
	"fmt"
	"testing"

	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/Coding-Brownies/todo/internal/repo/dbrepo"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)

	task := &entity.Task{
		Description: "lel",
		Done:        true,
	}
	err = r.Add(task)
	assert.NoError(t, err)

	res, err := r.List()
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, res[0].Description, "lel")
	assert.Equal(t, res[0].Done, true)
}

func TestCheck(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)

	task := &entity.Task{
		Description: "lel",
		Done:        false,
	}
	err = r.Add(task)
	assert.NoError(t, err)

	res, err := r.List()
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	err = r.Check(&res[0])
	assert.NoError(t, err)

	res, err = r.List()
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, true, res[0].Done)
}

func TestEdit(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)

	task := &entity.Task{
		Description: "spesa",
		Done:        false,
	}
	err = r.Add(task)
	assert.NoError(t, err)

	res, err := r.List()
	assert.NoError(t, err)

	err = r.Edit(&res[0], "ghes")
	assert.NoError(t, err)

	res, err = r.List()
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, "ghes", res[0].Description)
}

func TestSwap(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)

	tasks := []entity.Task{
		{
			Description: "calamaroA",
			Done:        false,
		},
		{
			Description: "cotechinoB",
			Done:        true,
		},
	}

	for _, task := range tasks {
		err = r.Add(&task)
		assert.NoError(t, err)
	}

	res, err := r.List()
	assert.NoError(t, err)

	err = r.Swap(&res[0], &res[1])
	assert.NoError(t, err)

	res, err = r.List()
	assert.NoError(t, err)

	assert.Len(t, res, 2)
	assert.Equal(t, res[0].Description, "cotechinoB")
	assert.Equal(t, res[0].Done, true)
	assert.Equal(t, res[1].Description, "calamaroA")
	assert.Equal(t, res[1].Done, false)
}

func TestUndo(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)
	tasks := []entity.Task{
		{
			Description: "Ale",
			Done:        false,
		},
		{
			Description: "Burberone",
			Done:        true,
		},
	}
	for _, task := range tasks {
		err = r.Add(&task)
		assert.NoError(t, err)
	}
	res, err := r.List()
	assert.NoError(t, err)
	// effettuo una modifica: idA true
	err = r.Check(&res[0])
	assert.NoError(t, err)
	// effettuo il ctrl+z
	err = r.Undo()
	assert.NoError(t, err)
	// stampa della lista aggiornata dei task
	res, err = r.List()
	assert.NoError(t, err)
	// il task deve essere come prima dell'ultima modifica
	assert.Equal(t, false, res[0].Done)
	assert.Equal(t, "Ale", res[0].Description)

	// effettuo una modifica: swap fra i due task
	err = r.Swap(&res[0], &res[1])
	assert.NoError(t, err)
	// effettuo il ctrl+z
	err = r.Undo()
	assert.NoError(t, err)
	// stampa della lista aggiornata dei task
	res, err = r.List()
	assert.NoError(t, err)
	// il task deve essere come prima dell'ultima modifica
	assert.True(t, res[0].Position.Before(res[1].Position))
}

func GetTaskStatusHash(r *dbrepo.DBRepo) (string, error) {
	tasks, err := r.List()
	if err != nil {
		return "", err
	}
	s := ""
	for _, t := range tasks {
		s += fmt.Sprint(t.Done) + "-" + t.Description + "-" + t.Position.String() + "|"
	}
	return s, nil
	// res := sha256.Sum256([]byte(s))
	// return fmt.Sprintf("%x", res), nil
}

func TestStronza(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)

	var (
		task1 entity.Task
		task2 entity.Task
		task3 entity.Task
	)

	mods := []func(){
		func() { r.Add(&task1) },
		func() { r.Edit(&task1, "albero") },
		func() { r.Add(&task2) },
		func() { r.Edit(&task2, "bar") },
		func() { r.Add(&task3) },
		func() { r.Edit(&task3, "cane") },
		func() { r.Check(&task1) },
		func() { r.Swap(&task1, &task2) },
		func() { r.Swap(&task1, &task3) },
		func() { r.Check(&task1) },
		func() { r.Swap(&task1, &task3) },
	}

	statuses := make([]string, len(mods))
	for i := 0; i < len(mods); i++ {
		status, err := GetTaskStatusHash(r)
		assert.NoError(t, err)
		statuses[i] = status

		mods[i]()
	}

	for i := len(mods) - 1; i >= 0; i-- {
		err = r.Undo()
		assert.NoError(t, err)
		status, err := GetTaskStatusHash(r)
		assert.NoError(t, err)
		assert.Equal(t, statuses[i], status, i)
	}

	// stampa della lista aggiornata dei task
	res, err := r.List()
	assert.NoError(t, err)
	// il task deve essere come prima dell'ultima modifica
	assert.True(t, len(res) == 0)
}
