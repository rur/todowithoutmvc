package todonomvc

type Todo struct {
	Active  bool
	ID      string
	Content string
}

type Todos []Todo

func (t Todos) ClearCompleted() Todos {
	activeOnly := t[:0]
	for _, todo := range t {
		if todo.Active {
			activeOnly = append(activeOnly, todo)
		}
	}
	return activeOnly
}

type Repository map[string][]Todo

func NewMemoryRepo() Repository {
	return make(Repository)
}

func (r *Repository) NewTodos() Todos {
	return make(Todos, 0)
}
