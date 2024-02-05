package handlers

import (
	"html/template"
	"net/http"
)

type Film struct {
	Title    string
	Director string
}

func DisplayFilms(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("public/index.html"))
	films := map[string][]Film{
		"Films": {
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "Blade Runner", Director: "Ridley Scott"},
			{Title: "The Thing", Director: "John Carpenter"},
		},
	}
	tmpl.Execute(w, films)
}
