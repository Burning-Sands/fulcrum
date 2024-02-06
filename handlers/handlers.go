package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// type Film struct {
// 	Title    string
// 	Director string
// }

func DisplayFilms(w http.ResponseWriter, r *http.Request) {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	m := map[interface{}]interface{}{}

	tmpl := template.Must(template.ParseFiles("public/index.html"))
	err = yaml.Unmarshal([]byte(yamlFile), &m)
	if err != nil {
		log.Fatal("error")
	}
	tmpl.Execute(w, m)
}

// func DisplayFilms(w http.ResponseWriter, r *http.Request) {
// 	tmpl := template.Must(template.ParseFiles("public/index.html"))
// 	films := map[string][]Film{
// 		"Films": {
// 			{Title: "The Godfather", Director: "Francis Ford Coppola"},
// 			{Title: "Blade Runner", Director: "Ridley Scott"},
// 			{Title: "The Thing", Director: "John Carpenter"},
// 		},
// 	}
// 	tmpl.Execute(w, films)
// }
