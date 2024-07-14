package main

import (
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

func (a *application) handlerDisplayIndex() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := a.templateCache["index.html"]
		tmpl.ExecuteTemplate(w, "base", a.templateData)
	}

	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := a.templateCache["display-values.html"]
		tmpl.ExecuteTemplate(w, "display-values", a.templateData)
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayOptions() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		pathValue := r.PathValue("option")

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

		default:
			a.clientError(w, http.StatusBadRequest)
		}

		w.Header().Add("HX-Trigger", "valuesChanged")
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerApplyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		// TODO Implement encoding functionality for all yaml files

		v, err := encodeTemplateData(a.templateData.Values)
		if err != nil {
			a.serverError(w, r, err)
		}
		c, err := encodeTemplateData(a.templateData.Chart)
		if err != nil {
			a.serverError(w, r, err)
		}

		git, err := gitlab.NewClient(*a.gitlabToken)
		if err != nil {
			a.serverError(w, r, err)
		}

		valuesFileName := "values.yaml"
		chartFileName := "Chart.yaml"

		pid := "fulcrum29/argoapps"
		brName := "test-branch"

		cf := []*gitlab.CommitActionOptions{
			{
				Action:   gitlab.Ptr(gitlab.FileCreate),
				Content:  gitlab.Ptr(v),
				FilePath: gitlab.Ptr(valuesFileName),
			},
			{
				Action:   gitlab.Ptr(gitlab.FileCreate),
				Content:  gitlab.Ptr(c),
				FilePath: gitlab.Ptr(chartFileName),
			},
		}

		commitOpts := &gitlab.CreateCommitOptions{
			Branch:        gitlab.Ptr(brName),
			StartBranch:   gitlab.Ptr("master"),
			CommitMessage: gitlab.Ptr("Add service files"),
			Actions:       cf,
		}

		_, res, err := git.Commits.CreateCommit(pid, commitOpts)
		if err != nil {
			a.serverError(w, r, err)
			a.clientError(w, http.StatusBadRequest)
		}
		if res.StatusCode == http.StatusOK {
			w.WriteHeader(http.StatusOK)
			tmpl := a.templateCache["apply-values.html"]
			tmpl.Execute(w, res)
		}
	}
	return http.HandlerFunc(fn)
}
