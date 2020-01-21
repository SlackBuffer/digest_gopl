package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}

		color := r.Form["color"][0]
		w.Header().Set("Content-Type", "image/svg+xml")
		svg(w, color)
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// http://localhost:8000/?color=yellow
