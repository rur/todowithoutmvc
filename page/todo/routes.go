package todo

import (
	"github.com/rur/todonomvc/page"
	"github.com/rur/treetop"
)

func Routes(server page.Server, m page.Mux, renderer *treetop.Renderer) {
	pageView := renderer.NewView(
		"page/todo/templates/index.templ.html",
		server.Bind(todoPageHandler),
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

	m.Handle("/", todo.PartialHandler().Include(footer))
	m.Handle("/active", todo.PartialHandler().Include(footer))
	m.Handle("/completed", todo.PartialHandler().Include(footer))

	m.HandleFunc("/clear", clearHandler(server))
	m.HandleFunc("/create", createHandler(server))
	m.HandleFunc("/toggle", toggleHandler(server))
}
