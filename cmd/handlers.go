package main

import (
	"bytes"
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

func (a *application) handlerDisplayIndex() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Display index")

		tmpl := a.templateCache["index.html"]
		tmpl.ExecuteTemplate(w, "base", a.templateData)
	}

	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		a.logger.Info("Display values")

		tmpl := a.templateCache["display-values.html"]
		tmpl.ExecuteTemplate(w, "display-values", a.templateData)
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayOptions() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		pathValue := r.PathValue("option")
		a.logger.Info("Display options", "option", pathValue)

		tmpl := a.templateCache["service-options.html"]
		err := tmpl.ExecuteTemplate(w, pathValue, a.templateData)
		if err != nil {
			a.clientError(w, http.StatusBadRequest)
		}

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

		writer := new(bytes.Buffer)
		encoder := yaml.NewEncoder(writer)
		encoder.SetIndent(2)
		encoder.Encode(*a.templateData.Values)
		encoder.Close()
		v := writer.String()

		git, err := gitlab.NewClient(*a.gitlabToken)
		if err != nil {
			a.serverError(w, r, err)
		}

		fileName := "values-output.yaml"
		pid := "fulcrum29/argoapps"
		brName := "test-branch"
		brCf := &gitlab.CreateBranchOptions{
			Branch: gitlab.Ptr(brName),
			Ref:    gitlab.Ptr("master"),
		}

		cf := &gitlab.UpdateFileOptions{
			Branch:        gitlab.Ptr(brName),
			Content:       gitlab.Ptr(v),
			CommitMessage: gitlab.Ptr("Modify test file"),
		}

		_, _, err = git.Branches.CreateBranch(pid, brCf)
		if err != nil {
			a.serverError(w, r, err)
			a.clientError(w, http.StatusBadRequest)
		}
		_, _, err = git.RepositoryFiles.UpdateFile(pid, fileName, cf)
		if err != nil {
			a.serverError(w, r, err)
			a.clientError(w, http.StatusBadRequest)
		}
	}
	return http.HandlerFunc(fn)
}
