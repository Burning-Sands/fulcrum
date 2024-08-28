package main

import "net/http"

func (app *application) logRequests(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("Received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *application) restoreSessionTemplateData(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		if app.sessionManager.Exists(r.Context(), "templateData") {

			*app.templateData = app.sessionManager.Get(r.Context(), "templateData").(TemplateData)
			app.logger.Info("Restored session data")
		} else {
			app.templateData = &TemplateData{}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
