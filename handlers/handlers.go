package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	"gopkg.in/yaml.v3"
)

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
	Tolerations               []interface{} `yaml:"tolerations"`
	TopologySpreadConstraints []interface{} `yaml:"topologySpreadConstraints"`
	Ports                     []interface{} `yaml:"ports"`
	Service                   struct {
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
	PodDisruptionBudget struct {
		Enabled      bool `yaml:"enabled"`
		MinAvailable int  `yaml:"minAvailable"`
	} `yaml:"podDisruptionBudget"`
}

func NewValues() *Values {
	return &Values{}
}

func DisplayIndex(values *Values) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("public/index.html"))
		tmpl.ExecuteTemplate(w, "index", *values)
	}

	return http.HandlerFunc(fn)
}

func DisplayValues(values *Values) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("public/index.html"))
		tmpl.ExecuteTemplate(w, "display-values", &values)
		fmt.Println("Display", values)
	}
	return http.HandlerFunc(fn)
}

func (v *Values) ModifyValues() http.Handler {

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

func ApplyValues(values *Values) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		fileName := "values-output.yaml"
		writer, err := os.Create(fileName)
		if err != nil {
			panic("Unable to create the output file")
		}
		encoder := yaml.NewEncoder(writer)
		encoder.SetIndent(2)
		encoder.Encode(*values)
		encoder.Close()
	}
	return http.HandlerFunc(fn)
}
