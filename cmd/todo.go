// create a folder cmd with the main.go
package main

import (
	"fmt"
	"os"

	"github.com/Coding-Brownies/todo/config"
	"github.com/Coding-Brownies/todo/internal/app"
	"github.com/Coding-Brownies/todo/internal/repo/dbrepo"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	r, err := dbrepo.New(os.Getenv("HOME") + "/.local/share/store.db")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	a := app.New(cfg, r)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "live")
	}
	// i secondi sono le cose dopo add
	// tipo go run main.go add "ciao patata"
	err = a.Run(os.Args[1], os.Args[2:]...)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
