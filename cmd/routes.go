package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static"))

	// define handlers
	router.Handle("GET /static/", http.StripPrefix("/static", fs))
	router.Handle("GET /{$}", app.handlerDisplayIndex())
	router.Handle("POST /edit/{option}", app.handlerModifyValues())
	router.Handle("GET /display-values", app.handlerDisplayValues())
	router.Handle("GET /apply", app.handlerApplyValues())
	router.Handle("GET /service-options/{option}", app.handlerDisplayOptions())

	return router
}
