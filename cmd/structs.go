package main

type Values struct {
	NameOverride     string `yaml:"nameOverride"`
	FullnameOverride string `yaml:"fullnameOverride"`
	ReplicaCount     int    `yaml:"replicaCount"`
	Annotations      struct {
	} `yaml:"annotations"`
	PodAnnotations struct {
	} `yaml:"podAnnotations"`
	Image struct {
		Repository string `yaml:"repository"`
		Tag        string `yaml:"tag"`
		PullPolicy string `yaml:"pullPolicy"`
	} `yaml:"image"`
	Env       []interface{} `yaml:"env"`
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
	VolumeMounts interface{} `yaml:"volumeMounts"`
	Affinity     struct {
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
	Ports       []interface{} `yaml:"ports"`
	Service     struct {
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

func newValues() *Values {
	return &Values{}
}
