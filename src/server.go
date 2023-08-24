package main

import (
	"flag"
	"log"
	"net/http"

	handlers "github.com/PhilipBatten/go-htmx-websockets/src/handlers"
)

func main() {
	flag.Parse()
	h := newHub()
	router := http.NewServeMux()
	router.HandleFunc("/", handlers.HomeHandler)
	router.Handle("/ws", wsHandler{h: h})
	log.Fatal(http.ListenAndServe(":3000", router))
}
