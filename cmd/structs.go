package main

import "log/slog"

type application struct {
	logger      *slog.Logger
	values      *Values
	gitlabToken *string
}

func NewApplication(l *slog.Logger, v *Values, gt *string) *application {
	return &application{
		logger:      l,
		values:      v,
		gitlabToken: gt,
	}
}

type EnvVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}

type Values struct {
	NameOverride     string `yaml:"nameOverride,omitempty"`
	FullnameOverride string `yaml:"fullnameOverride,omitempty"`
	ReplicaCount     int    `yaml:"replicaCount,omitempty"`
	Annotations      struct {
	} `yaml:"annotations,omitempty"`
	PodAnnotations struct {
	} `yaml:"podAnnotations,omitempty"`
	Image struct {
		Repository string `yaml:"repository,omitempty"`
		Tag        string `yaml:"tag,omitempty"`
		PullPolicy string `yaml:"pullPolicy,omitempty"`
	} `yaml:"image"`
	Env       []EnvVariable `yaml:"env"`
	EnvFrom   []interface{} `yaml:"envFrom"`
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
			} `yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
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
	} `yaml:"metrics"`
	Ingress struct {
	} `yaml:"ingress"`
	Hpa struct {
		Enabled     bool          `yaml:"enabled"`
		MinReplicas int           `yaml:"minReplicas"`
		MaxReplicas int           `yaml:"maxReplicas"`
		Metrics     []interface{} `yaml:"metrics"`
	} `yaml:"hpa"`
}

func NewValues() *Values {
	return &Values{}
}
