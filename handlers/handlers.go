package handlers

import (
	"html/template"
	"net/http"
	"os"

	// "github.com/fulcrum29/fulcrum/yamleditor"
	"gopkg.in/yaml.v3"
)

type Values struct {
	Uhc struct {
		Replicas string `yaml:"replicas"`
		Image    struct {
			Repository string `yaml:"repository"`
			Tag        string `yaml:"tag"`
		} `yaml:"image"`
		Resources struct {
			Memory struct {
				Requests int `yaml:"requests"`
				Limits   int `yaml:"limits"`
			} `yaml:"memory"`
		} `yaml:"resources"`
		Hpa struct {
			Enabled  bool
			Replicas struct {
				Min int `yaml:"min"`
				Max int `yaml:"max"`
			} `yaml:"replicas"`
		} `yaml:"hpa"`
	} `yaml:"uhc"`
}

func NewValues() *Values {
  return &Values{}
}

var values Values


func DisplayIndex(w http.ResponseWriter, r *http.Request) {


	tmpl := template.Must(template.ParseFiles("public/index.html"))
	tmpl.ExecuteTemplate(w, "index", values)
}

func DisplayValues(w http.ResponseWriter, r *http.Request) {


	tmpl := template.Must(template.ParseFiles("public/index.html"))
	tmpl.ExecuteTemplate(w, "display-values", values)
}

func ModifyValues(w http.ResponseWriter, r *http.Request) {

	var (
		repository = &values.Uhc.Image.Repository
		tag        = &values.Uhc.Image.Tag
		replicas   = &values.Uhc.Replicas
	)
	*repository = r.PostFormValue("image")
	*tag = r.PostFormValue("tag")
	*replicas = r.PostFormValue("replicas")

	w.Header().Add("HX-Trigger", "valuesChanged")
}


func ApplyValues(w http.ResponseWriter, r *http.Request) {

  fileName := "values-output.yaml"

	writer, err := os.Create(fileName)
	if err != nil {
		panic("Unable to create the output file")
	}
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	encoder.Encode(values)
	encoder.Close()

}














