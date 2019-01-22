package app

type Todo struct {
	Active  bool
	ID      string
	Content string
}

type Todos interface {
	ActiveOnly() Todos
	ActiveCount() int
	CompletedCount() int
	List() []Todo
	CompletedOnly() Todos
	AddEntry(string) (Todos, error)
	GetEntry(string) (*Todo, bool)
	UpdateEntry(Todo) (Todos, error)
	RemoveEntry(string) (Todos, error)
}
