// create a folder entity with a file for each used entity (ex: task.go which is a struct)
package entity

const CheckToDo = "◻"
const CheckDone = "◼"

type Task struct {
	ID			string
	Done        bool
	Description string
}

func (t Task) FilterValue() string { return "" }
