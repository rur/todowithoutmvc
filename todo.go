package todonomvc

type Todo struct {
	Active  bool
	ID      string
	Content string
}

type Todos []Todo

func (t Todos) ActiveOnly() Todos {
	filtered := t[:0]
	for _, todo := range t {
		if todo.Active {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

func (t Todos) CompletedOnly() Todos {
	filtered := t[:0]
	for _, todo := range t {
		if !todo.Active {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

type Repository map[string][]Todo

func NewMemoryRepo() Repository {
	return make(Repository)
}

func (r *Repository) NewTodos() Todos {
	return make(Todos, 0)
}
