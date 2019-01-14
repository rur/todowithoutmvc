package todo

import (
	"github.com/rur/todowithoutmvc/page"
	"github.com/rur/treetop"
)

func Routes(cxt page.Context, m page.Mux, renderer *treetop.Renderer) {
	pageView := renderer.NewView(
		"page/todo/templates/index.templ.html",
		todoPageHandler,
	)
	footer := pageView.DefaultSubView(
		"footer",
		"page/todo/templates/footer.templ.html",
		cxt.Bind(footerHandler),
	)
	todo := pageView.DefaultSubView(
		"main",
		"page/todo/templates/todos.templ.html",
		cxt.Bind(todoHandler),
	)
	editItem := renderer.NewView(
		"page/todo/templates/edit.templ.html",
		cxt.Bind(editTodoHandler),
	)

	m.Handle("/", treetop.ViewHandler(todo, footer))
	m.Handle("/active", treetop.ViewHandler(todo, footer))
	m.Handle("/completed", treetop.ViewHandler(todo, footer))
	m.Handle("/edit", treetop.ViewHandler(editItem).FragmentOnly())
}
