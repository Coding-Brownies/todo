package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal"
	"github.com/Coding-Brownies/todo/internal/bubble"
	"github.com/Coding-Brownies/todo/internal/bubble/components/bin"
	"github.com/Coding-Brownies/todo/internal/bubble/components/task"
	"github.com/Coding-Brownies/todo/internal/entity"
	"github.com/charmbracelet/bubbles/key"
)

type App struct {
	repo internal.Repo
	cfg  *config.Config
}

// funzione new restituisce un riferimento ad App
func New(cfg *config.Config, r internal.Repo) *App {

	return &App{
		repo: r,
		cfg:  cfg,
	}
}

// metodo run
func (a *App) Run(cmd string, args ...string) error {

	if cmd == "ls" {
		if len(args) != 0 {
			return errors.New("list accepts no argument")
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
			return errors.New("add accepts only one argument")
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
		cfg := a.cfg

		if len(args) != 0 {
			return errors.New("live accepts no argument")
		}
		keyMap := bubble.NewKeyMap(cfg)

		b := bubble.New(
			cfg,
			a.repo,
			keyMap,
			task.NewModel(
				&task.KeyMap{
					Check: key.NewBinding(
						bubble.WithKeys(cfg.Check...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Check), "(un)check"),
					),
					SwapUp: key.NewBinding(
						bubble.WithKeys(cfg.SwapUp...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.SwapUp), "swap up"),
					),
					SwapDown: key.NewBinding(
						bubble.WithKeys(cfg.SwapDown...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.SwapDown), "swap down"),
					),
					Remove: key.NewBinding(
						bubble.WithKeys(cfg.Remove...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Remove), "remove"),
					),
					Insert: key.NewBinding(
						bubble.WithKeys(cfg.Insert...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Insert), "insert"),
					),
					Up: key.NewBinding(
						bubble.WithKeys(cfg.Up...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Up), "go up"),
					),
					Down: key.NewBinding(
						bubble.WithKeys(cfg.Down...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Down), "go down"),
					),

					Edit: key.NewBinding(
						bubble.WithKeys(cfg.Edit...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Edit), "edit"),
					),
					EditExit: key.NewBinding(
						bubble.WithKeys(cfg.EditExit...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.EditExit), "to exit"),
					),
				},
				a.repo,
			),
			bin.NewModel(
				&bin.KeyMap{
					Up: key.NewBinding(
						bubble.WithKeys(cfg.Up...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Up), "go up"),
					),
					Down: key.NewBinding(
						bubble.WithKeys(cfg.Down...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Down), "go down"),
					),
					Restore: key.NewBinding(
						bubble.WithKeys(cfg.Restore...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.Restore), "restore"),
					),
					EmptyBin: key.NewBinding(
						bubble.WithKeys(cfg.EmptyBin...),
						key.WithHelp(bubble.ReplaceSymbols(cfg.EmptyBin), "empty the bin"),
					),
				},
				a.repo,
			),
		)

		return b.Run()
	}

	return errors.New("command not found")
}
