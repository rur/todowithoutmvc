package page

import (
	"github.com/rur/treetop"
)

func Routes(cxt Context, m Mux, exec treetop.ViewExecutor) {
	pageView := treetop.NewView(
		"page/templates/index.templ.html",
		todoPageHandler,
	)
	footer := pageView.NewDefaultSubView(
		"footer",
		"page/templates/footer.templ.html",
		cxt.Bind(footerHandler),
	)
	todo := pageView.NewDefaultSubView(
		"main",
		"page/templates/todos.templ.html",
		cxt.Bind(todoHandler),
	)
	editItem := treetop.NewView(
		"page/templates/edit.templ.html",
		cxt.Bind(editTodoHandler),
	)

	m.Handle("/", exec.NewViewHandler(todo, footer))
	m.Handle("/active", exec.NewViewHandler(todo, footer))
	m.Handle("/completed", exec.NewViewHandler(todo, footer))
	m.Handle("/edit", exec.NewViewHandler(editItem).FragmentOnly())
}
