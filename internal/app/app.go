package app

import (
	"errors"
	"fmt"

	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/entity"
)

type App struct {
	repo Repo
}

// funzione new restituisce un riferimento ad App
func New(r Repo) *App {
	return &App{
		repo: r,
	}
}

// metodo run
func (a *App) Run(cmd string, args ...string) error {
	if cmd == "ls" {
		if len(args) != 0 {
			return errors.New("list accept no argument")
		}
		tasks, err := a.repo.List()
		if err != nil {
			return err
		}
		for i, t := range tasks {
			if t.Done {
				fmt.Print(entity.CheckDone)
			} else {
				fmt.Print(entity.CheckToDo)
			}
			fmt.Println(" ", i, t.Description)
		}

		return nil
	}
	if cmd == "add" {
		if len(args) != 1 {
			return errors.New("add accept only one argument")
		}
		t := entity.Task{
			Description: args[0],
			Done:        false,
		}
		err := a.repo.Add(
			&t,
		)
		if err != nil {
			return err
		}
		return nil
	}
	if cmd == "delete" {
		if len(args) != 1 {
			return errors.New("delete accept only one argument")
		}
		return a.repo.Delete(args[0])
	}
	if cmd == "check" {
		if len(args) != 1 {
			return errors.New("check accept only one argument")
		}
		return a.repo.Check(args[0])
	}
	if cmd == "uncheck" {
		if len(args) != 1 {
			return errors.New("uncheck accept only one argument")
		}
		return a.repo.Uncheck(args[0])
	}

	if cmd == "edit" {
		if len(args) != 2 {
			return errors.New("edit accept 2 arguments")
		}
		return a.repo.Edit(args[0], args[1])
	}

	if cmd == "live" {
		if len(args) != 0 {
			return errors.New("live accept no argument")
		}
		tasks, err := a.repo.List()
		if err != nil {
			return err
		}

		res := bubble.Run(tasks)
		return a.repo.Store(res)
	}

	return errors.New("command not found")
}
