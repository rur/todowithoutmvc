package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rur/todowithoutmvc/app"
	"github.com/rur/todowithoutmvc/page"
	"github.com/rur/treetop"
)

var (
	addr = ":8000"
)

func main() {
	mux := http.NewServeMux()

	// static files
	js := http.FileServer(http.Dir("js"))
	mux.Handle("/js/", http.StripPrefix("/js/", js))
	css := http.FileServer(http.Dir("css"))
	mux.Handle("/css/", http.StripPrefix("/css/", css))
	modules := http.FileServer(http.Dir("node_modules"))
	mux.Handle("/node_modules/", http.StripPrefix("/node_modules/", modules))

	// server maintains all state and configuration
	server := app.NewServer(app.NewMemoryRepo())

	// Treetop view config
	page.Routes(
		page.NewContext(server),
		mux,
		treetop.NewRenderer(treetop.DefaultTemplateExec),
	)

	// CRUD handlers
	// POST/Redirect/GET is used for all side-effect endpoints
	mux.HandleFunc("/clear", app.ClearHandler(server))
	mux.HandleFunc("/create", app.CreateHandler(server))
	mux.HandleFunc("/toggle", app.ToggleHandler(server))
	mux.HandleFunc("/toggle-all", app.ToggleAllHandler(server))
	mux.HandleFunc("/update", app.UpdateHandler(server))
	mux.HandleFunc("/remove", app.RemoveHandler(server))

	fmt.Printf("Starting github.com/rur/todowithoutmvc server at %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
