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
		Page           string
		ActiveCount    int
		CompletedCount int
		Label          string
	}{
		ActiveCount:    rsc.Todos.ActiveCount(),
		CompletedCount: rsc.Todos.CompletedCount(),
	}
	if data.ActiveCount == 1 {
		data.Label = "item left"
	} else {
		data.Label = "items left"
	}
	if strings.HasPrefix(req.RequestURI, "/active") {
		data.Page = "active"
	} else if strings.HasPrefix(req.RequestURI, "/completed") {
		data.Page = "completed"
	} else {
		data.Page = "all"
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
		Todos: filtered.List(),
	}
}

// Doc: Purge all non active todos and redirect afterwards
func clearHandler(server page.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
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
		redirect := req.Referer()
		if redirect == "" {
			redirect = "/"
		}
		http.Redirect(w, req, redirect, http.StatusSeeOther)
	}
}

// Doc: Create a new todo entry
func createHandler(server page.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := server.LoadTodos(req)
		if key == "" {
			key = page.CreateTodoCookie(w)
		}
		value := strings.TrimSpace(req.FormValue("todo"))
		updated, err := todos.AddEntry(value)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating todo entry: %s", err.Error()), http.StatusBadRequest)
			return
		}
		if err := server.SaveTodos(key, updated); err != nil {
			http.Error(w, "Error saving todo list", http.StatusInternalServerError)
			return
		}
		redirect := req.Referer()
		if redirect == "" {
			redirect = "/"
		}
		http.Redirect(w, req, redirect, http.StatusSeeOther)
	}
}

// Doc: Update existing todo entry
func toggleHandler(server page.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := server.LoadTodos(req)
		if key == "" {
			http.Error(w, "Todo list was not found", http.StatusBadRequest)
			return
		}
		itemID := strings.TrimSpace(req.URL.Query().Get("item"))

		todo, ok := todos.GetEntry(itemID)
		if !ok {
			http.Error(w, fmt.Sprintf("Entry was not found for ID: %s", itemID), http.StatusBadRequest)
			return
		}

		if strings.ToLower(strings.TrimSpace(req.FormValue("completed"))) != "completed" {
			todo.Active = true
		} else {
			todo.Active = false
		}
		if updated, err := todos.UpdateEntry(*todo); err != nil {
			http.Error(w, "Error saving todo list", http.StatusInternalServerError)
			return
		} else if err := server.SaveTodos(key, updated); err != nil {
			http.Error(w, "Error saving todo list", http.StatusInternalServerError)
			return
		}

		redirect := req.Referer()
		if redirect == "" {
			redirect = "/"
		}
		http.Redirect(w, req, redirect, http.StatusSeeOther)
	}
}
