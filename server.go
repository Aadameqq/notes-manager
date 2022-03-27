package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world")
	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
