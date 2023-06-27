package internal

import "github.com/Coding-Brownies/todo/internal/entity"

type Repo interface {
	List() ([]entity.Task, error)
	Add(t *entity.Task) error
	Delete(ID string) error
	Check(ID string) error
	Uncheck(ID string) error
	Edit(ID string, newDescription string) error
	Swap(IDa, IDb string) error
}
