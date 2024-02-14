package main

import (
	"fmt"
	"os"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

func main() {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	var yamlNode yaml.Node
	var candidateNode yqlib.CandidateNode
	yaml.Unmarshal([]byte(yamlFile), &yamlNode)
	fmt.Println(candidateNode.UnmarshalYAML(yamlNode.Content[0], make(map[string]*yqlib.CandidateNode)))
	
	fmt.Println(candidateNode.GetKey())

	// define handlers
	// http.HandleFunc("/", handlers.DisplayNodes)
	// http.HandleFunc("/add-service/", handlers.AddService)

	// log.Fatal(http.ListenAndServe(":8080", nil))
}
