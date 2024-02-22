package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	var yamlNode yaml.Node
	// appendingYamlNode := yaml.Node{
	// 	Kind: 4, Style: 0, Value: "test"}
	//var candidateNode yqlib.CandidateNode
	yaml.Unmarshal([]byte(yamlFile), &yamlNode)

	//fmt.Println(candidateNode.UnmarshalYAML(yamlNode.Content[0], make(map[string]*yqlib.CandidateNode)))

	// fmt.Println(candidateNode.GetKey())

	// define handlers
	// http.HandleFunc("/", handlers.DisplayNodes)
	// http.HandleFunc("/add-service/", handlers.AddService)

	// log.Fatal(http.ListenAndServe(":8080", nil))
	//yamlNode.Content[0].Content = append(yamlNode.Content[0].Content, &appendingYamlNode)
	image := iterateOverYamlNode(&yamlNode)
	fmt.Println(image)

}

func iterateOverYamlNode(node *yaml.Node) *yaml.Node {
	var s *yaml.Node
	for _, v := range node.Content {
		switch v.Kind {
		case yaml.SequenceNode:
			//fmt.Println("Sequence Node", v.Kind)
			iterateOverYamlNode(v)
		case yaml.MappingNode:
			//fmt.Println("Mapping Node", v.Kind)
			iterateOverYamlNode(v)
		case yaml.ScalarNode:
			fmt.Println(v)
			if v.Value == "image" {
				s = v
			} else {
				fmt.Println(v.Value)
			}
			
		}

	}
	return s
}
