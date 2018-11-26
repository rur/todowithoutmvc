package todo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rur/todonomvc"
	"github.com/rur/todonomvc/page"
	"github.com/rur/treetop"
)

// todo (page)
// Doc: Todo *No* MVC base template
func todoPageHandler(rsc page.Resources, rsp treetop.Response, req *http.Request) interface{} {
	return struct {
		Footer interface{}
		Main   interface{}
	}{
		Footer: rsp.HandlePartial("footer", req),
		Main:   rsp.HandlePartial("main", req),
	}
}

// footer (default partial)
// Extends: footer
// Doc: Status and controls for todo list
func footerHandler(rsc page.Resources, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		Page    string
		Message string
	}{}
	if strings.HasPrefix(req.RequestURI, "/active") {
		data.Page = "active"
	} else if strings.HasPrefix(req.RequestURI, "/completed") {
		data.Page = "completed"
	} else {
		data.Page = "all"
	}
	left := 0
	for _, t := range rsc.Todos {
		if t.Active {
			left = left + 1
		}
	}
	if left == 1 {
		data.Message = "1 item left"
	} else {
		data.Message = fmt.Sprintf("%d items left", left)
	}

	return data
}

// handler for all todo list GET requests
// Extends: main
// Doc: List of all todos, filter based upon path
func todoHandler(rsc page.Resources, rsp treetop.Response, req *http.Request) interface{} {
	filtered := rsc.Todos[:0]
	if strings.HasPrefix(req.RequestURI, "/active") {
		for _, todo := range rsc.Todos {
			if todo.Active {
				filtered = append(filtered, todo)
			}
		}
	} else if strings.HasPrefix(req.RequestURI, "/completed") {
		for _, todo := range rsc.Todos {
			if !todo.Active {
				filtered = append(filtered, todo)
			}
		}
	} else {
		filtered = rsc.Todos
	}

	return struct {
		Todos []todonomvc.Todo
	}{
		Todos: filtered,
	}
}

// Doc: Purge all non active todos and redirect afterwards
func clearHandler(server page.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		todos, key := server.LoadTodos(req)
		if key == "" {
			http.Error(w, "Todo list was not found", http.StatusBadRequest)
			return
		}
		activeOnly := todos.ClearCompleted()
		if err := server.SaveTodos(key, activeOnly); err != nil {
			http.Error(w, "Error saving todo list", http.StatusInternalServerError)
			return
		}
		redirect := req.URL.Query().Get("redirect")
		if redirect == "" {
			redirect = "/"
		}
		http.Redirect(w, req, redirect, http.StatusSeeOther)
	}
}
