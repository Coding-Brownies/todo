package mock

import (
	"github.com/Coding-Brownies/todo/app"
	"github.com/Coding-Brownies/todo/entity"
)

var _ app.Repo = &MockRepo{}

type MockRepo struct{}

func New() *MockRepo {
	return &MockRepo{}
}
func (m *MockRepo) List() []entity.Task {
	return []entity.Task{
		{
			Description: "marameo",
		},
	}
}

func (m *MockRepo) Add(e *entity.Task) error {
	return nil
}

//remove
