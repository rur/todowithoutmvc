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

	// footer
	footer := pageView.DefaultSubView(
		"footer",
		"page/todo/templates/footer/footer.templ.html",
		cxt.Bind(footerHandler),
	)

	// main
	todo := pageView.DefaultSubView(
		"main",
		"page/todo/templates/main/todos.templ.html",
		cxt.Bind(todoHandler),
	)

	edit := renderer.NewView(
		"page/todo/templates/edit.templ.html",
		cxt.Bind(editTodoHandler),
	)

	m.Handle("/", treetop.ViewHandler(todo, footer))
	m.Handle("/active", treetop.ViewHandler(todo, footer))
	m.Handle("/completed", treetop.ViewHandler(todo, footer))
	m.Handle("/edit", treetop.ViewHandler(edit).FragmentOnly())
}
