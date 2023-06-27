package mock

import (
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/entity"
)

var _ internal.Repo = &MockRepo{}

type MockRepo struct{}

func New() *MockRepo {
	return &MockRepo{}
}

func (m *MockRepo) List() ([]entity.Task, error) {
	return []entity.Task{
		{
			Description: "marameo",
		},
	}, nil
}

func (m *MockRepo) Add(t *entity.Task) error {
	return nil
}

func (m *MockRepo) Delete(id string) error {
	return nil
}

func (m *MockRepo) Check(id string) error {
	return nil
}

func (m *MockRepo) Uncheck(id string) error {
	return nil
}

func (m *MockRepo) Edit(id string, newDescription string) error {
	return nil
}

// Swap implements app.Repo
func (*MockRepo) Swap(IDa string, IDb string) error {
	return nil
}
