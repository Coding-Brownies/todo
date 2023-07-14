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
	r, err := dbrepo.New(":memory:", 10)
	assert.NoError(t, err)

	task := &entity.Task{
		Description: "lel",
		Done:        true,
	}
	err = r.Add(task)
	assert.NoError(t, err)

	assert.Equal(t, task.Description, "lel")
	assert.Equal(t, task.Done, true)
}

func TestCheck(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:", 10)
	assert.NoError(t, err)

	task := &entity.Task{
		Description: "lel",
		Done:        false,
	}
	err = r.Add(task)
	assert.NoError(t, err)

	err = r.Check(task)
	assert.NoError(t, err)

	assert.Equal(t, true, task.Done)
}

func TestEdit(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:", 10)
	assert.NoError(t, err)

	task := &entity.Task{
		Description: "spesa",
		Done:        false,
	}
	err = r.Add(task)
	assert.NoError(t, err)

	err = r.Edit(task, "ghes")
	assert.NoError(t, err)

	assert.Equal(t, "ghes", task.Description)
}

func TestSwap(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:", 10)
	assert.NoError(t, err)

	tasks := []*entity.Task{
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
		err = r.Add(task)
		assert.NoError(t, err)
	}

	err = r.Swap(tasks[0], tasks[1])
	assert.NoError(t, err)

	assert.True(t, tasks[1].Position.Before(tasks[0].Position))

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
}

func TestUndo(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:", -1)
	assert.NoError(t, err)

	var (
		task1 entity.Task
		task2 entity.Task
		task3 entity.Task
	)

	mods := []func(){
		func() { r.Add(&task1) },
		func() { r.Edit(&task1, "Ale") },
		func() { r.Add(&task2) },
		func() { r.Edit(&task2, "Burberone") },
		func() { r.Add(&task3) },
		func() { r.Edit(&task3, "Calamaro") },
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

func TestHistoryLen(t *testing.T) {
	t.Parallel()

	const historyLen = 2
	r, err := dbrepo.New(":memory:", historyLen)
	assert.NoError(t, err)

	var (
		task1 entity.Task
		task2 entity.Task
	)

	mods := []func(){
		func() { r.Add(&task1) },
		func() { r.Edit(&task1, "Ale") },
		func() { r.Add(&task2) },
		func() { r.Edit(&task2, "Burberone") },
	}

	statuses := make([]string, len(mods))
	for i := 0; i < len(mods); i++ {
		status, err := GetTaskStatusHash(r)
		assert.NoError(t, err)
		statuses[i] = status

		mods[i]()
	}

	for i := len(mods) - 1; i >= 0; i-- {
		_ = r.Undo()
		status, err := GetTaskStatusHash(r)
		assert.NoError(t, err)
		index := i
		if index < historyLen {
			index = historyLen
		}
		assert.Equal(t, statuses[index], status, i)
	}
}

func TestBin(t *testing.T) {
	t.Parallel()

	r, err := dbrepo.New(":memory:", 2)
	assert.NoError(t, err)

	task := entity.Task{
		Description: "ciaoo",
	}

	// create a task
	err = r.Add(&task)
	assert.Nil(t, err)

	// delete a task
	err = r.Delete(&task)
	assert.Nil(t, err)

	// now there shold be a task in the bin
	res, err := r.ListBin()
	assert.Nil(t, err)
	assert.Len(t, res, 1)

	// and no task in the tasks
	res, err = r.List()
	assert.Nil(t, err)
	assert.Len(t, res, 0)

	// empty the bin
	err = r.EmptyBin()
	assert.Nil(t, err)

	// now the bin should be empty
	res, err = r.ListBin()
	assert.Nil(t, err)
	assert.Len(t, res, 0)
}
