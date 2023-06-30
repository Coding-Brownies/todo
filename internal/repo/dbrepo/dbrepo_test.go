package dbrepo_test

import (
	"testing"

	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/Coding-Brownies/todo/internal/repo/dbrepo"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
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

	err = r.Delete(res[0].ID)
	assert.NoError(t, err)
}

func TestCheck(t *testing.T) {
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

	err = r.Check(res[0].ID)
	assert.NoError(t, err)

	res, err = r.List()
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, true, res[0].Done)

	err = r.Delete(res[0].ID)
	assert.NoError(t, err)
}

func TestEdit(t *testing.T) {
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

	err = r.Edit(res[0].ID, "ghes")
	assert.NoError(t, err)

	res, err = r.List()
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, "ghes", res[0].Description)

	err = r.Delete(res[0].ID)
	assert.NoError(t, err)
}

func TestSwap(t *testing.T) {
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

	err = r.Delete(res[0].ID)
	assert.NoError(t, err)
	err = r.Delete(res[1].ID)
	assert.NoError(t, err)
}

func TestUndo(t *testing.T) {
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
	err = r.Check(res[0].ID)
	assert.NoError(t, err)
	// stampa della lista aggiornata dei task
	res, err = r.List()
	assert.NoError(t, err)
	// controlo che effettui la check
	assert.Equal(t, true, res[0].Done)
	assert.Equal(t, "Ale", res[0].Description)
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
	// stampa della lista aggiornata dei task
	res, err = r.List()
	assert.NoError(t, err)
	// controllo che abbia fatto la swap
	assert.Equal(t, "Burberone", res[0].Description)
	assert.Equal(t, "Ale", res[1].Description)
	// effettuo il ctrl+z
	err = r.Undo()
	assert.NoError(t, err)
	// stampa della lista aggiornata dei task
	res, err = r.List()
	assert.NoError(t, err)
	// il task deve essere come prima dell'ultima modifica
	assert.Equal(t, "Ale", res[0].Description)
	assert.Equal(t, "Burberone", res[1].Description)

	err = r.Delete(res[0].ID)
	assert.NoError(t, err)
	err = r.Delete(res[1].ID)
	assert.NoError(t, err)
}
