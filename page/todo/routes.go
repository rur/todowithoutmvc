package todo

import (
	"github.com/rur/todowithoutmvc/page"
	"github.com/rur/treetop"
)

func Routes(server page.Server, m page.Mux, renderer *treetop.Renderer) {
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

	// None treetop view handlers
	// I'm using POST redirect GET for all side-effect endpoints
	m.HandleFunc("/clear", clearHandler(server))
	m.HandleFunc("/create", createHandler(server))
	m.HandleFunc("/toggle", toggleHandler(server))
	m.HandleFunc("/toggle-all", toggleAllHandler(server))
}
