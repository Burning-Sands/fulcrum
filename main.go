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

  values := handlers.Values{}
  di := handlers.DisplayIndex(values)
  // fmt.Print(values)
  // fs := http.FileServer(http.Dir("./view"))
	// define handlers
  router := http.NewServeMux()
	router.Handle("/", di) 
	router.Handle("/display-values/", handlers.DisplayValues(values))
  router.Handle("/edit/", values.ModifyValues())
  router.Handle("/apply/", handlers.ApplyValues(values))
	log.Fatal(http.ListenAndServe(":8080", router))

}
