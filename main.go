package main

import (
	"log"
	"net/http"
	"github.com/fulcrum29/fulcrum/handlers"
)


func main() {
  
  values := handlers.Values{}
  // fs := http.FileServer(http.Dir("./view"))
	// define handlers
  router := http.NewServeMux()
	router.Handle("/", handlers.DisplayIndex(&values)) 
  router.Handle("/edit/", values.ModifyValues())
	router.Handle("/display-values/", handlers.DisplayValues(&values))
  router.Handle("/apply/", handlers.ApplyValues(&values))
	log.Fatal(http.ListenAndServe(":8080", router))

}
