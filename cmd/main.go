package main

import (
	"log"
	"net/http"
)


func main() {
  
  values := NewValues()

  router := http.NewServeMux()
  fs := http.FileServer(http.Dir("./ui/static"))

	// define handlers
  router.Handle("GET /static/", http.StripPrefix("/static", fs))
	router.Handle("GET /{$}", DisplayIndex(values)) 
  router.Handle("POST /edit", values.ModifyValues())
	router.Handle("GET /display-values", DisplayValues(values))
  router.Handle("GET /apply", ApplyValues(values))
	log.Fatal(http.ListenAndServe(":8080", router))

}
