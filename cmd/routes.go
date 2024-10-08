package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static"))

	// define handlers
	router.Handle("GET /static/", http.StripPrefix("/static", fs))
	router.Handle("GET /{$}", app.sessionManager.LoadAndSave(app.restoreSessionTemplateData(app.handlerDisplayIndex())))
	router.Handle("POST /edit/{option}", app.sessionManager.LoadAndSave(app.restoreSessionTemplateData(app.handlerModifyValues())))
	router.Handle("GET /display-values", app.sessionManager.LoadAndSave(app.restoreSessionTemplateData(app.handlerDisplayValues())))
	router.Handle("GET /apply", app.handlerApplyValues())
	router.Handle("GET /service-options/{option}", app.sessionManager.LoadAndSave(app.restoreSessionTemplateData(app.handlerDisplayOptions())))

	return app.logRequests(router)
}
