package main

import (
<<<<<<< HEAD
	"flag"
	"log"
=======
	"log/slog"
>>>>>>> 4c2aa1f (feat: add structred logging)
	"net/http"
  "os"
)

<<<<<<< HEAD
func main() {
=======

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

func main() {

  values := NewValues()
  logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

  app := application{
    logger: logger,
  }
>>>>>>> 4c2aa1f (feat: add structred logging)

	gitlabToken := flag.String("token", "", "Gitlab token for repo")
	flag.Parse()

	values := NewValues()

	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./ui/static"))

	// define handlers
<<<<<<< HEAD
	router.Handle("GET /static/", http.StripPrefix("/static", fs))
	router.Handle("GET /{$}", DisplayIndex(values))
	router.Handle("POST /edit", values.ModifyValues())
	router.Handle("GET /display-values", DisplayValues(values))
	router.Handle("GET /apply", ApplyValues(values, gitlabToken))
	log.Fatal(http.ListenAndServe(":8080", router))
=======
  router.Handle("GET /static/", http.StripPrefix("/static", fs))
	router.Handle("GET /{$}", app.DisplayIndex(values)) 
  router.Handle("POST /edit", app.ModifyValues(values))
	router.Handle("GET /display-values", app.DisplayValues(values))
  router.Handle("GET /apply", app.ApplyValues(values))


  // Start main handler (server)
  logger.Info("Starting server")
  err := http.ListenAndServe(":8080", router)

  logger.Error(err.Error())
  os.Exit(1)
>>>>>>> 4c2aa1f (feat: add structred logging)

}
