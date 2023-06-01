package jsonrepo

import "github.com/Coding-Brownies/todo/app"

var _ app.Repo = &JSONRepo{}

type JSONRepo struct {
	path string
}

func New(p string) *JSONRepo {
	return &JSONRepo{
		path: p,
	}
}
