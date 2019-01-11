package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rur/todowithoutmvc"
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

	server := todowithoutmvc.NewServer(
		todowithoutmvc.NewMemoryRepo(),
	)

	renderer := treetop.NewRenderer(treetop.DefaultTemplateExec)
	todo.Routes(server, m, renderer)

	// I'm using POST redirect GET for all side-effect endpoints
	m.HandleFunc("/clear", todowithoutmvc.ClearHandler(server))
	m.HandleFunc("/create", todowithoutmvc.CreateHandler(server))
	m.HandleFunc("/toggle", todowithoutmvc.ToggleHandler(server))
	m.HandleFunc("/toggle-all", todowithoutmvc.ToggleAllHandler(server))

	fmt.Printf("Starting github.com/rur/todowithoutmvc server at %s", addr)
	// Bind to an addr and pass our router in
	log.Fatal(http.ListenAndServe(addr, m))
}
