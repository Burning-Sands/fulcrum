package main

import (
	"log"
	"net/http"

	"github.com/fulcrum29/fulcrum/handlers"
)

func main() {

	// define handlers
	http.HandleFunc("/", handlers.DisplayNodes)
	http.HandleFunc("/add-service/", handlers.AddService)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// f, err := os.OpenFile("values.yaml", os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	panic(err)
	// }

}
