package todo

import (
	"github.com/rur/todowithoutmvc"
	"github.com/rur/todowithoutmvc/page"
	"github.com/rur/treetop"
)

func Routes(server todowithoutmvc.Server, m page.Mux, renderer *treetop.Renderer) {
	pageView := renderer.NewView(
		"page/todo/templates/index.templ.html",
		todoPageHandler,
	)

	// footer
	footer := pageView.DefaultSubView(
		"footer",
		"page/todo/templates/footer/footer.templ.html",
		server.Bind(footerHandler),
	)

	// main
	todo := pageView.DefaultSubView(
		"main",
		"page/todo/templates/main/todos.templ.html",
		server.Bind(todoHandler),
	)

	m.Handle("/", treetop.ViewHandler(todo, footer))
	m.Handle("/active", treetop.ViewHandler(todo, footer))
	m.Handle("/completed", treetop.ViewHandler(todo, footer))
}
