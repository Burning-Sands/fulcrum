package main

import (
	"log"
	"net/http"

	"github.com/fulcrum29/fulcrum/handlers"
)

func main() {

	// yamlFile, err := os.ReadFile("values.yaml")
	// if err != nil {
	// 	panic(err)
	// }

	// m := map[interface{}]interface{}{}
	// films := yaml.Unmarshal([]byte(yamlFile), &m)
	// if films != nil {
	// 	log.Fatal("error")
	// }
	// fmt.Print(m)

	// define handlers
	http.HandleFunc("/", handlers.DisplayFilms)
	// //http.HandleFunc("/add-film/", handlers.AddFilm)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
