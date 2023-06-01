package app

import (
	"errors"
	"fmt"

	"github.com/Coding-Brownies/todo/entity"
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
		tasks := a.repo.List()
		fmt.Println(tasks)
		return nil
	}
	if cmd == "add" {
		var t entity.Task
		err := a.repo.Add(
			&t,
		)
		return err
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

	return nil
}
