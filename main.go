package main

import (
	"log"
	"net/http"

	"github.com/fulcrum29/fulcrum/handlers"
)

func main() {
  // values := handlers.NewValues()
	// playground
	// var valuesFile yamleditor.Values
	// file, err := os.ReadFile("values.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// yaml.Unmarshal([]byte(file), &valuesFile)
	// fmt.Printf("%+v\n", valuesFile)

	// define handlers
  router := http.NewServeMux() 
	router.HandleFunc("/", handlers.DisplayNodes)
  router.HandleFunc("/edit/", handlers.ModifyValues)
  router.HandleFunc("/apply/", handlers.ApplyValues)
	log.Fatal(http.ListenAndServe(":8080", router))

}
