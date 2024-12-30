package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	"github.com/fulcrum29/fulcrum/pkg/templatedata"
	github "github.com/google/go-github/github"
	"gopkg.in/yaml.v3"
)

func (a *application) handlerDisplayIndex() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := a.htmlTemplateCache["index.html"]
		err := tmpl.ExecuteTemplate(w, "base", a.templateData)
		if err != nil {
			a.serverError(w, r, err)
		}
	}

	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := a.htmlTemplateCache["display-values.html"]
		err := tmpl.ExecuteTemplate(w, "display-values", a.templateData)
		if err != nil {
			a.serverError(w, r, err)
		}
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerDisplayOptions() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		pathValue := r.PathValue("option")

		tmpl := a.htmlTemplateCache["service-options.html"]
		err := tmpl.ExecuteTemplate(w, pathValue, a.templateData)
		if err != nil {
			a.serverError(w, r, err)
		}

	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerModifyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		var (
			chartName          = &a.templateData.Chart.Name
			uhc                = &a.templateData.Values.Uhc
			repository         = &a.templateData.Values.Uhc.Image.Repository
			tag                = &a.templateData.Values.Uhc.Image.Tag
			replicas           = &a.templateData.Values.Uhc.ReplicaCount
			limits             = &a.templateData.Values.Uhc.Resources.Limits
			requests           = &a.templateData.Values.Uhc.Resources.Requests
			ports              = &a.templateData.Values.Uhc.Ports
			hpa                = &a.templateData.Values.Uhc.Hpa
			env                = &a.templateData.Values.Uhc.Env
			k8sRepo            = &a.templateData.K8sRepo
			serviceName        = &a.templateData.githubTemplate.Input.ServiceName
			serviceNameRuleStg = &a.templateData.githubTemplate.Input.ServiceNameRuleStg
			serviceNameRulePrd = &a.templateData.githubTemplate.Input.ServiceNameRulePrd
		)

		err := r.ParseForm()
		if err != nil {
			a.serverError(w, r, err)
			return
		}

		formGet := r.PostForm.Get
		pathValue := r.PathValue("option")

		switch pathValue {

		case "basic":
			*chartName = r.PostForm.Get("serviceName")
			*serviceName = r.PostForm.Get("serviceName")
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
				f, err := os.ReadFile("templates/nodeAffinitySpot.yaml")
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
				err = yaml.Unmarshal(f, uhc)
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
			} else if aff == "regular" {
				f, err := os.ReadFile("templates/nodeAffinityRegular.yaml")
				if err != nil {
					a.serverError(w, r, err)
					os.Exit(1)
				}
				err = yaml.Unmarshal(f, uhc)
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
			e := templatedata.EnvVariable{
				Name:  formGet("envName"),
				Value: formGet("envValue"),
			}
			*env = append(*env, e)

		case "github":
			*k8sRepo = r.PostForm.Get("k8sRepo")
			*serviceName = r.PostForm.Get("serviceName")
			*serviceNameRuleStg = r.PostForm.Get("serviceNameRuleStg")
			*serviceNameRulePrd = r.PostForm.Get("serviceNameRulePrd")

		default:
			a.clientError(w, errors.New("Wrong path, option doesn't exist"), 400)
		}

		a.sessionManager.Put(r.Context(), "templateData", *a.templateData)

		w.Header().Add("HX-Trigger", "valuesChanged")
	}
	return http.HandlerFunc(fn)
}

func (a *application) handlerApplyValues() http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		v, err := a.templateData.EncodeValues()
		if err != nil {
			a.serverError(w, r, err)
		}
		c, err := a.templateData.EncodeChart()
		if err != nil {
			a.serverError(w, r, err)
		}
		// at, err := a.populateArgoAppTemplate()
		// if err != nil {
		// 	a.serverError(w, r, err)
		// }
		// gt, err := a.templateData.EncodegithubTemplate()
		if err != nil {
			a.serverError(w, r, err)
		}

		var (
			serviceName    = a.templateData.Chart.Name
			valuesFilePath = fmt.Sprintf("services/%s/values.yaml", serviceName)
			chartFilePath  = fmt.Sprintf("services/%s/Chart.yaml", serviceName)
			// argoAppFilePath  = fmt.Sprintf("argocdapps/templates/%s.yaml", serviceName)
			pid    = "fulcrum29/argoapps"
			brName = fmt.Sprintf("Deployment_of_%s", serviceName)
		)

		cf := []*github.CommitActionOptions{
			{
				Action:   github.Ptr(github.FileCreate),
				Content:  github.Ptr(v),
				FilePath: github.Ptr(valuesFilePath),
			},
			{
				Action:   github.Ptr(github.FileCreate),
				Content:  github.Ptr(c),
				FilePath: github.Ptr(chartFilePath),
			},
			// {
			//
			// 	Action:   github.Ptr(github.FileCreate),
			// 	Content:  github.Ptr(at),
			// 	FilePath: github.Ptr(argoAppFilePath),
			// },
		}

		commitOpts := &github.CreateCommitOptions{
			Branch:        github.Ptr(brName),
			StartBranch:   github.Ptr("master"),
			CommitMessage: github.Ptr("Add service files"),
			Actions:       cf,
		}

		_, res, err := git.Commits.CreateCommit(pid, commitOpts)
		if err != nil {
			a.logger.Error(err.Error())
			tmpl := a.htmlTemplateCache["apply-values.html"]
			err := tmpl.ExecuteTemplate(w, "applyError", err.Error())
			if err != nil {
				a.serverError(w, r, err)
			}
		}
		a.logger.Info("Received response status from github", "Response", res.Status)
		if res.StatusCode == http.StatusCreated {
			a.templateData = templatedata.NewTemplateData()
			tmpl := a.htmlTemplateCache["apply-values.html"]
			err := tmpl.ExecuteTemplate(w, "applySuccess", nil)
			if err != nil {
				a.serverError(w, r, err)
			}

		}
	}
	return http.HandlerFunc(fn)
}
