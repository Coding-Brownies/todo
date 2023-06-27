package dbrepo_test

import (
	"testing"

	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/Coding-Brownies/todo/internal/repo/dbrepo"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	r, err := dbrepo.New("/tmp/store.db")
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
	r, err := dbrepo.New("/tmp/store.db")
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
	r, err := dbrepo.New("/tmp/store.db")
	assert.NoError(t, err)

	task := &entity.Task{
		ID:          "lollo",
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
	r, err := dbrepo.New("/tmp/store.db")
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

	err = r.Swap(res[0].ID, res[1].ID)
	assert.NoError(t, err)

	res, err = r.List()
	assert.NoError(t, err)

	assert.Len(t, res, 2)
	assert.Equal(t, res[0].Description, "cotechinoB")

	err = r.Delete(res[0].ID)
	assert.NoError(t, err)
	err = r.Delete(res[1].ID)
	assert.NoError(t, err)
}
