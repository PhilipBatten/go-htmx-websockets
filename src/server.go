package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpl := template.Must(template.ParseFiles("resources/views/index.html"))
		tpl.Execute(w, r)
	})
}

func main() {
	flag.Parse()
	h := newHub()
	router := http.NewServeMux()
	router.Handle("/", HomeHandler())
	router.Handle("/ws", wsHandler{h: h})
	log.Fatal(http.ListenAndServe(":3000", router))
}
