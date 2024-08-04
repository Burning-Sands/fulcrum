package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	method := r.Method
	uri := r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, err error, status int) {

	http.Error(w, err.Error(), status)

}

func (app *application) populateArgoAppTemplate() (string, error) {

	serviceName := app.templateData.Chart.Name
	k8sRepo := app.templateData.K8sRepo
	t, err := os.ReadFile("templates/argocdAppTemplate.yaml")

	if err != nil {
		return "", err
	}

	t = bytes.ReplaceAll(t, []byte("replacemetemplate"), []byte(serviceName))
	t = bytes.ReplaceAll(t, []byte("replacemeproject"), []byte(k8sRepo))

	var data bytes.Buffer
	data.Write(t)
	at := data.String()
	return at, nil

}
func (app *application) populateGitlabCiTemplate() (string, error) {

	serviceName := app.templateData.Chart.Name
	t, err := os.ReadFile("templates/gitlabCiTemplate.yaml")

	if err != nil {
		return "", err
	}

	t = bytes.ReplaceAll(t, []byte("replacemetemplate"), []byte(serviceName))

	var data bytes.Buffer
	data.Write(t)
	gt := data.String()
	return gt, nil

}

// TODO Implement render helper function
// func (app *application) renderTemplate()

func newHtmlTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
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
			continue
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

func (app *application) encodeChart() (string, error) {

	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	defer encoder.Close()

	encoder.SetIndent(2)
	err := encoder.Encode(app.templateData.Chart)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (app *application) encodeValues() (string, error) {

	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	defer encoder.Close()

	encoder.SetIndent(2)
	err := encoder.Encode(app.templateData.Values)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
