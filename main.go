package main

import (
	"fmt"
	"log"
	"os"

	// git "github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing/transport/http"
	// "github.com/xanzy/go-gitlab"
	// "gopkg.in/yaml.v3"
	yptr "github.com/vmware-labs/yaml-jsonpointer"
	yaml "gopkg.in/yaml.v3"
)

func main() {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	keyPath := "/uhc/image"
	var yamlNode yaml.Node
	yaml.Unmarshal([]byte(yamlFile), &yamlNode)

	r, _ := yptr.Find(&yamlNode, keyPath)
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)

	r, _ = yptr.Find(&yamlNode, `/uhc/tag`)
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)
	r.SetString("1234")
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)

	f, err := os.Create("values.yaml")
	if err != nil {
		log.Fatalf("Problem opening file: %v", err)
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	encoder.Encode(yamlNode.Content[0])
	encoder.Close()

}
