package main

import (
	"log"
	"net/http"

	"github.com/fulcrum29/fulcrum/handlers"
)

func main() {
	// playground
	// var valuesFile yamleditor.Values
	// file, err := os.ReadFile("values.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// yaml.Unmarshal([]byte(file), &valuesFile)
	// fmt.Printf("%+v\n", valuesFile)

	// define handlers
	http.HandleFunc("/", handlers.DisplayNodes)
	http.HandleFunc("/add-service/", handlers.AddService)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
