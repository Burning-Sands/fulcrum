package main

import (
	"html/template"
	"log/slog"
)

type application struct {
	logger        *slog.Logger
	templateData  *TemplateData
	gitlabToken   *string
	templateCache map[string]*template.Template
}

func NewApplication(
	l *slog.Logger,
	td *TemplateData,
	gt *string,
	tc map[string]*template.Template,
) *application {
	return &application{
		logger:        l,
		templateData:  td,
		gitlabToken:   gt,
		templateCache: tc,
	}
}

type EnvVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}

type TemplateData struct {
	Chart  *Chart
	Values *Values
}

// func NewTemplateData() *TemplateData {
// 	return &TemplateData{}
// }

type Chart struct {
	Name         string            `yaml:"name"`
	AppVersion   string            `yaml:"appVersion"`
	Dependencies map[string]string `yaml:"dependencies"`
}

func NewChart() *Chart {
	deps := map[string]string{
		"uhc": "0.30.1",
	}
	return &Chart{
		AppVersion:   "0.1.0",
		Dependencies: deps,
	}
}

type Values struct {
	Uhc struct {
		NameOverride     string `yaml:"nameOverride,omitempty"`
		FullnameOverride string `yaml:"fullnameOverride,omitempty"`
		ReplicaCount     int    `yaml:"replicaCount,omitempty"`
		Annotations      struct {
		} `yaml:"annotations,omitempty"`
		PodAnnotations struct {
		} `yaml:"podAnnotations,omitempty"`
		Image struct {
			Repository string `yaml:"repository"`
			Tag        string `yaml:"tag"`
			PullPolicy string `yaml:"pullPolicy,omitempty"`
		} `yaml:"image"`
		Env       []EnvVariable `yaml:"env,omitempty"`
		EnvFrom   []interface{} `yaml:"envFrom,omitempty"`
		Resources struct {
			Limits struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"limits"`
			Requests struct {
				CPU    string `yaml:"cpu"`
				Memory string `yaml:"memory"`
			} `yaml:"requests"`
		} `yaml:"resources"`
		Affinity struct {
			NodeAffinity struct {
				RequiredDuringSchedulingIgnoredDuringExecution struct {
					NodeSelectorTerms []struct {
						MatchExpressions []struct {
							Key      string   `yaml:"key"`
							Operator string   `yaml:"operator"`
							Values   []string `yaml:"values"`
						} `yaml:"matchExpressions"`
					} `yaml:"nodeSelectorTerms"`
				} `yaml:"requiredDuringSchedulingIgnoredDuringExecution"`
			} `yaml:"nodeAffinity"`
		} `yaml:"affinity"`
		Tolerations []interface{} `yaml:"tolerations"`
		Ports       struct {
			Name          string `yaml:"name"`
			ContainerPort int    `yaml:"containerPort"`
			Protocol      string `yaml:"protocol"`
		} `yaml:"ports"`
		Service struct {
			Ports []struct {
				Name       string `yaml:"name"`
				TargetPort string `yaml:"targetPort"`
				Protocol   string `yaml:"protocol"`
				Port       int    `yaml:"port"`
			} `yaml:"ports"`
		} `yaml:"service"`
		Metrics struct {
			Enabled        bool `yaml:"enabled"`
			ServiceMonitor struct {
				Enabled bool `yaml:"enabled"`
			} `yaml:"serviceMonitor"`
		} `yaml:"metrics,omitempty"`
		Ingress struct {
		} `yaml:"ingress,omitempty"`
		Hpa struct {
			Enabled     bool          `yaml:"enabled"`
			MinReplicas int           `yaml:"minReplicas"`
			MaxReplicas int           `yaml:"maxReplicas"`
			Metrics     []interface{} `yaml:"metrics"`
		} `yaml:"hpa,omitempty"`
	} `yaml:"uhc"`
}

func NewValues() *Values {
	return &Values{}
}
