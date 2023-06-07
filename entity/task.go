// create a folder entity with a file for each used entity (ex: task.go which is a struct)
package entity

type Task struct{
	ID string `json:"-"`
	Done bool
	Description string
}