package dbrepo_test

import (
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

func TestStronza(t *testing.T) {
	t.Parallel()
	r, err := dbrepo.New(":memory:")
	assert.NoError(t, err)

	task1 := entity.Task{
		Description: "",
		Done:        false,
	}
	task2 := entity.Task{
		Description: "",
		Done:        false,
	}
	task3 := entity.Task{
		Description: "",
		Done:        false,
	}
	// modifica 1
	err = r.Add(&task1)
	assert.NoError(t, err)
	// modifica 2
	err = r.Edit(&task1, "albero")
	assert.NoError(t, err)
	// modifica 3
	err = r.Add(&task2)
	assert.NoError(t, err)
	// modifica 4
	err = r.Edit(&task2, "bar")
	assert.NoError(t, err)
	// modifica 5
	err = r.Add(&task3)
	assert.NoError(t, err)
	// modifica 6
	err = r.Edit(&task3, "cane")
	assert.NoError(t, err)
	// modifica 7
	err = r.Check(&task1)
	assert.NoError(t, err)
	// modifica 8
	err = r.Swap(&task1, &task2)
	assert.NoError(t, err)
	// modifica 9
	err = r.Swap(&task1, &task3)
	assert.NoError(t, err)
	// modifica 10
	err = r.Uncheck(&task1)
	assert.NoError(t, err)
	// modifica 11
	err = r.Swap(&task1, &task3)
	assert.NoError(t, err)

	// effettuo il ctrl+z
	for i := 0; i < 11; i++ {
		err = r.Undo()
		assert.NoError(t, err)
	}
	// stampa della lista aggiornata dei task
	res, err := r.List()
	assert.NoError(t, err)
	// il task deve essere come prima dell'ultima modifica
	assert.True(t, len(res) == 0)
}
