package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shorty API"))
	})
	log.Println("API Running on http://localhost:3000")

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
