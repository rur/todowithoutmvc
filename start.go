package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rur/todowithoutmvc/app"
	"github.com/rur/todowithoutmvc/page"
	"github.com/rur/todowithoutmvc/page/todo"
	"github.com/rur/treetop"
)

var (
	addr = ":8000"
)

func main() {
	m := http.NewServeMux()

	// static files
	js := http.FileServer(http.Dir("js"))
	m.Handle("/js/", http.StripPrefix("/js/", js))
	css := http.FileServer(http.Dir("css"))
	m.Handle("/css/", http.StripPrefix("/css/", css))
	modules := http.FileServer(http.Dir("node_modules"))
	m.Handle("/node_modules/", http.StripPrefix("/node_modules/", modules))

	// maintains all server state and configuration
	s := app.NewServer(app.NewMemoryRepo())

	todo.Routes(
		page.NewContext(s),
		m,
		treetop.NewRenderer(treetop.DefaultTemplateExec),
	)

	// I'm using POST redirect GET for all side-effect endpoints
	m.HandleFunc("/clear", app.ClearHandler(s))
	m.HandleFunc("/create", app.CreateHandler(s))
	m.HandleFunc("/toggle", app.ToggleHandler(s))
	m.HandleFunc("/toggle-all", app.ToggleAllHandler(s))
	m.HandleFunc("/update", app.UpdateHandler(s))

	fmt.Printf("Starting github.com/rur/todowithoutmvc server at %s", addr)
	// Bind to an addr and pass our router in
	log.Fatal(http.ListenAndServe(addr, m))
}
