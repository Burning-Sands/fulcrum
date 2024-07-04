package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

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
		tmpl.ExecuteTemplate(w, "base", a.templateData)
	}

	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("ui/html/pages/display-values.html"))
		tmpl.ExecuteTemplate(w, "display-values", a.templateData)
		a.logger.Info("Display values")
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayOptions() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		pathValue := r.PathValue("option")
		tmpl := template.Must(template.ParseFiles("ui/html/pages/service-options.html"))
		err := tmpl.ExecuteTemplate(w, pathValue, a.templateData)
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
			repository = &a.templateData.Values.Uhc.Image.Repository
			tag        = &a.templateData.Values.Uhc.Image.Tag
			replicas   = &a.templateData.Values.Uhc.ReplicaCount
			limits     = &a.templateData.Values.Uhc.Resources.Limits
			requests   = &a.templateData.Values.Uhc.Resources.Requests
			ports      = &a.templateData.Values.Uhc.Ports
			affinity   = &a.templateData.Values.Uhc.Affinity
			hpa        = &a.templateData.Values.Uhc.Hpa
			env        = &a.templateData.Values.Uhc.Env
		)

		err := r.ParseForm()
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
			return
		}

		formGet := r.PostForm.Get
		pathValue := r.PathValue("option")

		switch pathValue {

		case "basic":
			ports.ContainerPort, _ = strconv.Atoi(formGet("port-number"))
			*repository = r.PostForm.Get("repository")
			*tag = r.PostForm.Get("tag")
			*replicas, _ = strconv.Atoi(formGet("replicas"))
		case "resources":
			limits.CPU = formGet("cpu-limits")
			limits.Memory = formGet("memory-limits")
			requests.CPU = formGet("cpu-requests")
			requests.Memory = formGet("memory-requests")
		case "affinity-hpa":

			aff := formGet("affinity")

			if aff == "spot" {
				f, err := os.ReadFile("nodeAffinitySpot.yaml")
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
				err = yaml.Unmarshal(f, &affinity)
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
			} else if aff == "regular" {
				f, err := os.ReadFile("nodeAffinityRegular.yaml")
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
				err = yaml.Unmarshal(f, &affinity)
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
			} else {
				a.serverError(w, r, err)
				os.Exit(1)
			}

			if formGet("hpaEnabled") == "enabled" {
				hpa.Enabled = true
				a.templateData.Values.Uhc.ReplicaCount = 0
			} else {
				hpa.Enabled = false
			}

			hpa.MinReplicas, _ = strconv.Atoi(formGet("hpaMinReplicas"))
			hpa.MaxReplicas, _ = strconv.Atoi(formGet("hpaMaxReplicas"))

		case "env":
			e := EnvVariable{
				Name:  formGet("envName"),
				Value: formGet("envValue"),
			}
			*env = append(*env, e)
			a.logger.Info("Env variables", "Env", *env)

		default:
			a.clientError(w, http.StatusBadRequest)
		}

		a.logger.Info("Modify values", "path", pathValue)
		w.Header().Add("HX-Trigger", "valuesChanged")
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerApplyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		fileName := "values-output.yaml"
		//
		// writer, err := os.Create(fileName)
		//
		// if err != nil {
		// 	a.serverError(w, r, err)
		// 	os.Exit(1)
		// }
		//
		var writer bytes.Buffer
		encoder := yaml.NewEncoder(&writer)
		encoder.SetIndent(2)
		encoder.Encode(*a.templateData)
		encoder.Close()
		//
		// file, _ := os.ReadFile(fileName)
		// fileAsString := string(file)

		// out, err := yaml.Marshal(a.values)
		// if err != nil {
		// 	a.serverError(w, r, err)
		// }
		v := writer.String()

		git, err := gitlab.NewClient(*a.gitlabToken)
		if err != nil {
			a.serverError(w, r, err)
		}

		cf := &gitlab.UpdateFileOptions{
			Branch:        gitlab.Ptr("master"),
			Content:       gitlab.Ptr(v),
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
