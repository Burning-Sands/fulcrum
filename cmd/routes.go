package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static"))

	// define handlers
	router.Handle("GET /static/", http.StripPrefix("/static", fs))
	router.Handle("GET /{$}", app.DisplayIndex())
	router.Handle("POST /edit", app.ModifyValues())
	router.Handle("GET /display-values", app.DisplayValues())
	router.Handle("GET /apply", app.ApplyValues())
	router.Handle("GET /service-options/{option}", app.handlerDisplayOptions())

	return router
}
