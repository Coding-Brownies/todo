// create a folder cmd with the main.go
package main

import (
	"fmt"
	"os"

	"github.com/Coding-Brownies/todo/internal/app"
	"github.com/Coding-Brownies/todo/internal/repo/jsonrepo"
)

func main() {
	/* una volta finita dbrepo:
	r, err := dbrepo.New(os.Getenv("HOME") + "/.local/share/store.db")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	*/
	r := jsonrepo.New(os.Getenv("HOME") + "/.local/share/store.json")

	a := app.New(r)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "live")
	}
	// i secondi sono le cose dopo add
	// tipo go run main.go add "ciao patata"
	err := a.Run(os.Args[1], os.Args[2:]...)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
