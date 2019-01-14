package app

import (
	"fmt"
	"net/http"
	"strings"
)

// Doc: Purge all non active todos and redirect
func ClearHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := s.LoadTodos(req)
		if key == "" {
			http.Error(w, "Todo list was not found", http.StatusBadRequest)
			return
		}
		activeOnly := todos.ActiveOnly()
		if err := s.SaveTodos(key, activeOnly); err != nil {
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

// Doc: Create a new todo entry and redirect
func CreateHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := s.LoadTodos(req)
		if key == "" {
			key = CreateTodoCookie(w)
		}
		value := strings.TrimSpace(req.FormValue("todo"))
		updated, err := todos.AddEntry(value)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating todo entry: %s", err.Error()), http.StatusBadRequest)
			return
		}
		if err := s.SaveTodos(key, updated); err != nil {
			http.Error(w, "Error saving todo list", http.StatusInternalServerError)
			return
		}
		redirect := req.Referer()
		if redirect == "" || strings.Contains(redirect, "/completed") {
			redirect = "/"
		}
		http.Redirect(w, req, redirect, http.StatusSeeOther)
	}
}

// Doc: Toggle completeness of existing todo entry and redirect
func ToggleHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := s.LoadTodos(req)
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
		} else if err := s.SaveTodos(key, updated); err != nil {
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

// Doc: Toggle completeness of all todo entries and redirect
func ToggleAllHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := s.LoadTodos(req)
		if key == "" {
			http.Error(w, "Todo list was not found", http.StatusBadRequest)
			return
		}

		setComplete := strings.ToLower(strings.TrimSpace(req.FormValue("completed"))) == "completed"

		var err error
		for _, item := range todos.List() {
			item.Active = !setComplete
			todos, err = todos.UpdateEntry(item)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to update todo item #%s", item.ID), http.StatusInternalServerError)
				return
			}
		}

		if err := s.SaveTodos(key, todos); err != nil {
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

// Doc: update contents of a given todo item
func UpdateHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.Method) != "post" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		todos, key := s.LoadTodos(req)
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

		value := strings.TrimSpace(req.FormValue("content"))
		todo.Content = value

		if updated, err := todos.UpdateEntry(*todo); err != nil {
			http.Error(w, fmt.Sprintf("Error saving todo entry %s", err), http.StatusInternalServerError)
			return
		} else if err := s.SaveTodos(key, updated); err != nil {
			http.Error(w, fmt.Sprintf("Error saving todo list %s", err), http.StatusInternalServerError)
			return
		}

		redirect := req.Referer()
		if redirect == "" {
			redirect = "/"
		}
		http.Redirect(w, req, redirect, http.StatusSeeOther)
	}
}
