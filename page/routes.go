package page

import (
	"github.com/rur/treetop"
)

func Routes(cxt Context, m Mux, renderer *treetop.Renderer) {
	pageView := renderer.NewView(
		"page/templates/index.templ.html",
		todoPageHandler,
	)
	footer := pageView.DefaultSubView(
		"footer",
		"page/templates/footer.templ.html",
		cxt.Bind(footerHandler),
	)
	todo := pageView.DefaultSubView(
		"main",
		"page/templates/todos.templ.html",
		cxt.Bind(todoHandler),
	)
	editItem := renderer.NewView(
		"page/templates/edit.templ.html",
		cxt.Bind(editTodoHandler),
	)

	m.Handle("/", treetop.ViewHandler(todo, footer))
	m.Handle("/active", treetop.ViewHandler(todo, footer))
	m.Handle("/completed", treetop.ViewHandler(todo, footer))
	m.Handle("/edit", treetop.ViewHandler(editItem).FragmentOnly())
}
