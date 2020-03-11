package page

import (
	"net/http"

	"github.com/rur/todowithoutmvc/app"
	"github.com/rur/treetop"
)

// Interface for router configuration
type Mux interface {
	Handle(pattern string, handler http.Handler)
}

type Context interface {
	Bind(TodoHandler) treetop.ViewHandlerFunc
}

type cxt struct {
	srv app.Server
}

func NewContext(s app.Server) Context {
	return &cxt{s}
}

type TodoHandler func(app.Todos, treetop.Response, *http.Request) interface{}

func (c *cxt) Bind(f TodoHandler) treetop.ViewHandlerFunc {
	return func(rsp treetop.Response, req *http.Request) interface{} {
		// load user todo list from repo based upon request cookies
		// pass to handler.
		//
		// Note that these handlers have no way to making changes
		// to the todo list.
		todos, _ := c.srv.LoadTodos(req)
		return f(todos, rsp, req)
	}
}
