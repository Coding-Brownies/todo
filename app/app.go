// create a folder app that has a struct App which containes a IRepo interface and has some functions like list, add etc

package app

import "fmt"

type App struct{
	repo Repo
}

// funzione new restituisce un riferimento a
func New(r Repo) *App {
	return &App{
		repo: r,
	}
}

// metodo run
func (a *App) Run(cmd string) {
	if cmd == "ls"{
		tasks := a.repo.List()
		fmt.Println(tasks)
		return
	}
}
