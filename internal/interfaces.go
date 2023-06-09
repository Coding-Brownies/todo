package internal

import "github.com/Coding-Brownies/todo/internal/entity"

type Repo interface {
	List() ([]entity.Task, error)
	Add(t *entity.Task) error
	Delete(t *entity.Task) error
	Check(t *entity.Task) error
	Edit(t *entity.Task, newDescription string) error
	Swap(taskA, taskB *entity.Task) error
	Undo() error
	ListBin() ([]entity.Task, error)
	Restore(task *entity.Task) error
	EmptyBin() error
}
