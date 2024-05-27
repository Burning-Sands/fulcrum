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
  logger *slog.Logger
  values *Values
  gitlabToken *string
} 

func (app *application) DisplayIndex() http.Handler {

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
		tmpl.ExecuteTemplate(w, "base", *app.values)
	}

	return http.HandlerFunc(fn)
}

func (app *application) DisplayValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("ui/html/pages/display-values.html"))
		tmpl.ExecuteTemplate(w, "display-values", *app.values)
    app.logger.Info("Display values")
	}
	return http.HandlerFunc(fn)
}

func (app *application) ModifyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var (
			repository = &app.values.Image.Repository
			tag        = &app.values.Image.Tag
			replicas   = &app.values.ReplicaCount
		)
		*repository = r.PostFormValue("image")
		*tag = r.PostFormValue("tag")
		*replicas, _ = strconv.Atoi(r.PostFormValue("replicas"))

    app.logger.Info("Modify values")
		w.Header().Add("HX-Trigger", "valuesChanged")
	}
	return http.HandlerFunc(fn)
}

func (app *application) ApplyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		fileName := "values-output.yaml"

		writer, err := os.Create(fileName)

		if err != nil {
			app.serverError(w, r, err)
      os.Exit(1)
		}

		encoder := yaml.NewEncoder(writer)
		encoder.SetIndent(2)
		encoder.Encode(*app.values)
		encoder.Close()

		file, _ := os.ReadFile(fileName)
		fileAsString := string(file)

		git, err := gitlab.NewClient(*app.gitlabToken)
    if err != nil {
			app.serverError(w, r, err)
    }

		cf := &gitlab.UpdateFileOptions{
			Branch:        gitlab.Ptr("master"),
			Content:       gitlab.Ptr(fileAsString),
			CommitMessage: gitlab.Ptr("Adding a test file"),
		}

    _, _, err = git.RepositoryFiles.UpdateFile("fulcrum29/argoapps", fileName, cf)
	  if err != nil {
			app.serverError(w, r, err)
	  }
	}
	return http.HandlerFunc(fn)
}
