package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {

	http.Error(w, http.StatusText(status), status)

}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(".ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		if name == "index.html" {
			files := []string{
				"./ui/html/base.html",
			}
			files = append(files, pages...)
			tmpl, err := template.ParseFiles(files...)

			if err != nil {
				return nil, err
			}

			cache[name] = tmpl
		}

		files := []string{
			page,
		}
		tmpl, err := template.ParseFiles(files...)

		if err != nil {
			return nil, err
		}

		cache[name] = tmpl

	}
	return cache, nil
}
