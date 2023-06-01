// create a folder cmd with the main.go
package main

import (
	"fmt"
	"os"

	"github.com/Coding-Brownies/todo/app"
	"github.com/Coding-Brownies/todo/repo/mock"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Wrong format")
		return
	}
	r := mock.New()
	a := app.New(r)

	// i secondi sono le cose dopo add
	// tipo go run main.go add "ciao patata"
	err := a.Run(os.Args[1], os.Args[1:]...)

	if err != nil {
		fmt.Println("Error: ", err)
	}
}
