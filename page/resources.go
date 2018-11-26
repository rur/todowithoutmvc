package page

import (
	"net/http"

	"github.com/rur/todowithoutmvc"
	"github.com/rur/treetop"
)

type Resources struct {
	Todos todowithoutmvc.Todos
}

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type ResourcesHandler func(Resources, treetop.Response, *http.Request) interface{}
