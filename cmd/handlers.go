package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

type application struct {
	logger      *slog.Logger
	values      *Values
	gitlabToken *string
}

func (a *application) handlerDisplayIndex() http.Handler {

	// Declare templated files
	templateFiles := []string{
		"ui/html/base.html",
		"ui/html/pages/index.html",
		"ui/html/pages/apply-values.html",
		"ui/html/pages/service-options.html",
		"ui/html/pages/display-values.html",
	}

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles(templateFiles...))
		tmpl.ExecuteTemplate(w, "base", *a.values)
	}

	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("ui/html/pages/display-values.html"))
		tmpl.ExecuteTemplate(w, "display-values", *a.values)
		a.logger.Info("Display values")
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayOptions() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		pathValue := r.PathValue("option")
		tmpl := template.Must(template.ParseFiles("ui/html/pages/service-options.html"))
		err := tmpl.ExecuteTemplate(w, pathValue, *a.values)
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
		}
		a.logger.Info("Display options")

	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerModifyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var (
			repository = &a.values.Image.Repository
			tag        = &a.values.Image.Tag
			replicas   = &a.values.ReplicaCount
			limits     = &a.values.Resources.Limits
			requests   = &a.values.Resources.Requests
			ports      = &a.values.Ports
			// hpa        = &a.values.Hpa.Enabled
		)
		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		pathValue := r.PathValue("option")

		if pathValue == "basic" {
			ports.ContainerPort, _ = strconv.Atoi(r.PostForm.Get("port-number"))
			*repository = r.PostForm.Get("repository")
			*tag = r.PostForm.Get("tag")
			*replicas, _ = strconv.Atoi(r.PostForm.Get("replicas"))
		} else if pathValue == "resources" {
			limits.CPU = r.PostForm.Get("cpu-limits")
			limits.Memory = r.PostForm.Get("memory-limits")
			requests.CPU = r.PostForm.Get("cpu-requests")
			requests.Memory = r.PostForm.Get("memory-requests")
		} else {
			a.clientError(w, http.StatusBadRequest)
		}
		// *hpa = r.PostForm.Get("hpa-enabled")

		a.logger.Info("Modify values")
		w.Header().Add("HX-Trigger", "valuesChanged")
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerApplyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		fileName := "values-output.yaml"

		writer, err := os.Create(fileName)

		if err != nil {
			a.serverError(w, r, err)
			os.Exit(1)
		}

		encoder := yaml.NewEncoder(writer)
		encoder.SetIndent(2)
		encoder.Encode(*a.values)
		encoder.Close()

		file, _ := os.ReadFile(fileName)
		fileAsString := string(file)

		git, err := gitlab.NewClient(*a.gitlabToken)
		if err != nil {
			a.serverError(w, r, err)
		}

		cf := &gitlab.UpdateFileOptions{
			Branch:        gitlab.Ptr("master"),
			Content:       gitlab.Ptr(fileAsString),
			CommitMessage: gitlab.Ptr("Adding a test file"),
		}

		_, _, err = git.RepositoryFiles.UpdateFile("fulcrum29/argoapps", fileName, cf)
		if err != nil {
			a.serverError(w, r, err)
			a.clientError(w, http.StatusBadRequest)
		}
	}
	return http.HandlerFunc(fn)
}
