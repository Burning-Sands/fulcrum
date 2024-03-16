package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"

	"github.com/fulcrum29/fulcrum/yamleditor"
	"gopkg.in/yaml.v3"
)

var yamlEdit yamleditor.YamlOperator

func DisplayNodes(w http.ResponseWriter, r *http.Request) {

	if reflect.ValueOf(yamlEdit).IsZero() {
		file, err := os.ReadFile("values.yaml")
		if err != nil {
			panic(err)
		}
		yaml.Unmarshal([]byte(file), &yamlEdit.YamlNode)
	}

	yamlEdit.Buffer.Reset()
	yamleditor.IterateOverYamlNode(&yamlEdit.YamlNode, &yamlEdit.Buffer)

	tmpl := template.Must(template.ParseFiles("public/index.html"))
	tmpl.Execute(w, yamlEdit.Buffer.String())
}

func AddService(w http.ResponseWriter, r *http.Request) {

	yamleditor.IterateOverYamlNode(&yamlEdit.YamlNode, &yamlEdit.Buffer)
	fmt.Println(yamlEdit.Buffer.String())
	title := r.PostFormValue("title")
	director := r.PostFormValue("director")
	yamleditor.ChangeYamlNodeValue(&yamlEdit.YamlNode, title, director)
	w.Header().Add("HX-Trigger", "newValue")

	// tmpl := template.Must(template.ParseFiles("public/index.html"))
	// tmpl.ExecuteTemplate(w, "film-list-element", m)
}
