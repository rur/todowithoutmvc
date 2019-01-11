package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Repository map[string]Todos

func NewMemoryRepo() Repository {
	return make(Repository)
}

func (r *Repository) NewTodos() Todos {
	return &todos{
		list:   make([]Todo, 0),
		lastID: 0,
	}
}

type todos struct {
	list   []Todo
	lastID int
}

func (t *todos) ActiveOnly() Todos {
	filtered := []Todo{}
	for _, todo := range t.list {
		if todo.Active {
			filtered = append(filtered, todo)
		}
	}
	return &todos{
		list:   filtered,
		lastID: t.lastID,
	}
}

func (t *todos) ActiveCount() int {
	count := 0
	for _, todo := range t.list {
		if todo.Active {
			count = count + 1
		}
	}
	return count
}

func (t *todos) CompletedCount() int {
	count := 0
	for _, todo := range t.list {
		if !todo.Active {
			count = count + 1
		}
	}
	return count
}

func (t *todos) List() []Todo {
	return append([]Todo{}, t.list...)
}

func (t *todos) CompletedOnly() Todos {
	filtered := []Todo{}
	for _, todo := range t.list {
		if !todo.Active {
			filtered = append(filtered, todo)
		}
	}
	return &todos{
		list:   filtered,
		lastID: t.lastID,
	}
}

func (t *todos) AddEntry(msg string) (Todos, error) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return t, errors.New("Cannot create empty Todo item")
	} else {
		return &todos{
			list: append(t.list, Todo{
				Active:  true,
				Content: msg,
				ID:      strconv.FormatInt(int64(t.lastID+1), 10),
			}),
			lastID: t.lastID + 1,
		}, nil
	}
}

func (t *todos) GetEntry(id string) (*Todo, bool) {
	for _, t := range t.list {
		if t.ID == id {
			return &t, true
		}
	}
	return nil, false
}

func (t *todos) UpdateEntry(nue Todo) (Todos, error) {
	list := append([]Todo{}, t.list...)
	found := false
	for i, ti := range list {
		if ti.ID == nue.ID {
			list[i] = nue
			found = true
			break
		}
	}
	if !found {
		return t, fmt.Errorf("Entry not found with item ID %s", nue.ID)
	}
	return &todos{
		list:   list,
		lastID: t.lastID,
	}, nil
}
