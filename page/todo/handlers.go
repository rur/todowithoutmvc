package todo

import (
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
		Page  string
		Count int
		Label string
	}{}
	if strings.HasPrefix(req.RequestURI, "/active") {
		data.Page = "active"
	} else if strings.HasPrefix(req.RequestURI, "/completed") {
		data.Page = "completed"
	} else {
		data.Page = "all"
	}
	for _, t := range rsc.Todos {
		if t.Active {
			data.Count = data.Count + 1
		}
	}
	if data.Count == 1 {
		data.Label = "item left"
	} else {
		data.Label = "items left"
	}
	return data
}

// handler for all todo list GET requests
// Extends: main
// Doc: List of all todos, filter based upon path
func todoHandler(rsc page.Resources, rsp treetop.Response, req *http.Request) interface{} {
	filtered := rsc.Todos
	if strings.HasPrefix(req.RequestURI, "/active") {
		filtered = rsc.Todos.ActiveOnly()
	} else if strings.HasPrefix(req.RequestURI, "/completed") {
		filtered = rsc.Todos.CompletedOnly()
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
		activeOnly := todos.ActiveOnly()
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
