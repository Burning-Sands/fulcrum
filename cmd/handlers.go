package main

import (
	"fmt"
	"html/template"
<<<<<<< HEAD
	"log"
=======
	"log/slog"
>>>>>>> 4c2aa1f (feat: add structred logging)
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)


type application struct {
  logger *slog.Logger
} 


func (app *application) DisplayIndex(values *Values) http.Handler {

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
		tmpl.ExecuteTemplate(w, "base", *values)
	}

	return http.HandlerFunc(fn)
}

func (app *application) DisplayValues(values *Values) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("ui/html/pages/display-values.html"))
		tmpl.ExecuteTemplate(w, "display-values", &values)
		fmt.Println("Display", values)
	}
	return http.HandlerFunc(fn)
}

func (app *application) ModifyValues(v *Values) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var (
			repository = &v.Image.Repository
			tag        = &v.Image.Tag
			replicas   = &v.ReplicaCount
		)
		*repository = r.PostFormValue("image")
		*tag = r.PostFormValue("tag")
		*replicas, _ = strconv.Atoi(r.PostFormValue("replicas"))

		fmt.Println("Modify", v)
		w.Header().Add("HX-Trigger", "valuesChanged")
	}
	return http.HandlerFunc(fn)
}

<<<<<<< HEAD
func ApplyValues(values *Values, gitlabToken *string) http.Handler {
=======
func (app *application) ApplyValues(values *Values) http.Handler {
>>>>>>> 4c2aa1f (feat: add structred logging)

	fn := func(w http.ResponseWriter, r *http.Request) {

		fileName := "values-output.yaml"

		writer, err := os.Create(fileName)

		if err != nil {
			app.logger.Error("Unable to create the output file")
      os.Exit(1)
		}

		encoder := yaml.NewEncoder(writer)
		encoder.SetIndent(2)
		encoder.Encode(*values)
		encoder.Close()

		file, _ := os.ReadFile(fileName)
		fileAsString := string(file)

<<<<<<< HEAD
		git, err := gitlab.NewClient(*gitlabToken)
		if err != nil {
			log.Fatal(err)
		}
=======
    git, err := gitlab.NewClient("")
    if err != nil {
		  app.logger.Error(err.Error())
    }
>>>>>>> 4c2aa1f (feat: add structred logging)

		cf := &gitlab.UpdateFileOptions{
			Branch:        gitlab.Ptr("master"),
			Content:       gitlab.Ptr(fileAsString),
			CommitMessage: gitlab.Ptr("Adding a test file"),
		}

<<<<<<< HEAD
		_, _, err = git.RepositoryFiles.UpdateFile("fulcrum29/argoapps", fileName, cf)
		if err != nil {
			log.Print(err)
		}
=======
    _, _, err = git.RepositoryFiles.UpdateFile("fulcrum29/argoapps", fileName, cf)
	  if err != nil {
		  app.logger.Error(err.Error())
	  }
>>>>>>> 4c2aa1f (feat: add structred logging)
	}
	return http.HandlerFunc(fn)
}
