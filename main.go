package main

import (
	"log"
	"net/http"
)

func main() {
	go h.run()
	go initDBListener()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
