package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/entity"
)

type App struct {
	repo Repo
	cfg  *config.Config
}

// funzione new restituisce un riferimento ad App
func New(cfg *config.Config, r Repo) *App {

	return &App{
		repo: r,
		cfg:  cfg,
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
		for _, t := range tasks {
			if t.Done {
				fmt.Print(entity.CheckDone)
			} else {
				fmt.Print(entity.CheckToDo)
			}

			// remove multiple lines
			str := t.Description
			if idx := strings.Index(str, "\n"); idx != -1 {
				str = str[:idx] + "..."
			}
			fmt.Println(" ", str)
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

	if cmd == "live" {
		if len(args) != 0 {
			return errors.New("live accept no argument")
		}
		tasks, err := a.repo.List()
		if err != nil {
			return err
		}

		b := bubble.New(a.cfg)
		res := b.Run(tasks)
		return a.repo.Store(res)
	}

	return errors.New("command not found")
}
