// create a folder cmd with the main.go
package main

import (
	"fmt"
	"os"

	"github.com/Coding-Brownies/todo/app"
	"github.com/Coding-Brownies/todo/repo/jsonrepo"
)

func main() {

	r := jsonrepo.New(os.Getenv("HOME") + "/.local/share/store.json")
	a := app.New(r)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "ls")
	}
	// i secondi sono le cose dopo add
	// tipo go run main.go add "ciao patata"
	err := a.Run(os.Args[1], os.Args[2:]...)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
