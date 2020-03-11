package page

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rur/todowithoutmvc/app"
	"github.com/rur/treetop"
)

// todo (page)
// Doc: Todo *No* MVC base template
func todoPageHandler(rsp treetop.Response, req *http.Request) interface{} {
	return struct {
		Footer interface{}
		Main   interface{}
	}{
		Footer: rsp.HandleSubView("footer", req),
		Main:   rsp.HandleSubView("main", req),
	}
}

// footer (default partial)
// Extends: footer
// Doc: Status and controls for todo list
func footerHandler(todos app.Todos, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		Page           string
		ActiveCount    int
		CompletedCount int
		Label          string
	}{
		ActiveCount:    todos.ActiveCount(),
		CompletedCount: todos.CompletedCount(),
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
func todoHandler(todos app.Todos, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		Todos        []app.Todo
		AllCompleted bool
	}{}

	if strings.HasPrefix(req.RequestURI, "/active") {
		data.Todos = todos.ActiveOnly().List()
		data.AllCompleted = false
	} else if strings.HasPrefix(req.RequestURI, "/completed") {
		data.Todos = todos.CompletedOnly().List()
		data.AllCompleted = true
	} else {
		data.Todos = todos.List()
		data.AllCompleted = len(data.Todos) > 0 && len(data.Todos) == todos.CompletedCount()
	}

	// in the UI we want to show them in reverse order
	for i, j := 0, len(data.Todos)-1; i < j; i, j = i+1, j-1 {
		data.Todos[i], data.Todos[j] = data.Todos[j], data.Todos[i]
	}

	return data
}

// handler for loading an edit todo form
// Doc: configure a HTML form element which will be used to update the content of an existing todo
func editTodoHandler(todos app.Todos, rsp treetop.Response, req *http.Request) interface{} {
	todoID := req.URL.Query().Get("item")
	if todoID == "" {
		http.Error(rsp, "Missing Todo ID", http.StatusBadRequest)
	}

	if todo, found := todos.GetEntry(todoID); !found {
		// using `rsp` as a http.ResposeWriter will safely abort the Treetop response process.
		http.Error(rsp, fmt.Sprintf("Invalid todo ID: %s", todoID), http.StatusBadRequest)
		return nil
	} else {
		return todo
	}
}
