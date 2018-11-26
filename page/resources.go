package page

import (
	"net/http"

	"github.com/rur/todonomvc"
	"github.com/rur/treetop"
)

type Resources struct {
	Todos todonomvc.Todos
}

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type ResourcesHandler func(Resources, treetop.Response, *http.Request) interface{}
