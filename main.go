package main

import (
	"log"
	"net/http"
	"github.com/fulcrum29/fulcrum/handlers"
)

func main() {


	// define handlers
	http.HandleFunc("/", handlers.DisplayFilms)
	//http.HandleFunc("/add-film/", handlers.AddFilm)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
