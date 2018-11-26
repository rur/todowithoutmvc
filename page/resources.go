package page

import (
	"net/http"

	"github.com/rur/todonomvc"
	"github.com/rur/treetop"
)

type Resources struct {
	Todos  []todonomvc.Todo
	Update func(todonomvc.Todos) error
}

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)
}

type ResourcesHandler func(Resources, treetop.Response, *http.Request) interface{}