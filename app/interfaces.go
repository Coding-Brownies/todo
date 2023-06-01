package app

import "github.com/Coding-Brownies/todo/entity"

type Repo interface {
	List() []entity.Task
	Add(t *entity.Task) error
	Delete(ID string) error
	Check(ID string) error
	Uncheck(ID string) error
}
