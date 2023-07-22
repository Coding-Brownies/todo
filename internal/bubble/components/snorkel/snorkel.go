package snorkel

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var f *os.File

func Log(a ...any) {
	var err error

	f, err = tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, a...)
}
