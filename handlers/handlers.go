package handlers

import (
	"html/template"
	"net/http"

	"github.com/fulcrum29/fulcrum/yamleditor"
)

var yamlEdit yamleditor.YamlOperator

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

var values Values

func DisplayNodes(w http.ResponseWriter, r *http.Request) {

	// if reflect.ValueOf(values).IsZero() {
	// 	file, err := os.ReadFile("values.yaml")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	yaml.Unmarshal([]byte(file), &values)
	// }
	// fmt.Printf("%+v\n", values)

	tmpl := template.Must(template.ParseFiles("public/index.html"))
	tmpl.Execute(w, values)
}

func AddService(w http.ResponseWriter, r *http.Request) {

	var (
		repository = &values.Uhc.Image.Repository
		tag        = &values.Uhc.Image.Tag
		replicas   = &values.Uhc.Replicas
	)
	*repository = r.PostFormValue("image")
	*tag = r.PostFormValue("tag")
	*replicas = r.PostFormValue("replicas")

	w.Header().Add("HX-Trigger", "valuesChanged")
	// tmpl := template.Must(template.ParseFiles("public/index.html"))
	// tmpl.ExecuteTemplate(w, "film-list-element", m)
}
