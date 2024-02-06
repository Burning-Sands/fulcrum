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

func DisplayNodes(w http.ResponseWriter, r *http.Request) {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	m := map[string]interface{}{}

	tmpl := template.Must(template.ParseFiles("public/index.html"))
	err = yaml.Unmarshal([]byte(yamlFile), &m)
	if err != nil {
		log.Fatal("error")
	}
	tmpl.Execute(w, m)
}

func AddFilms(w http.ResponseWriter, r *http.Request) {

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
