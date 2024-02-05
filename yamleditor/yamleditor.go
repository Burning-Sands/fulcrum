package yamleditor

import (
	"fmt"
	"log"
	"os"

	yptr "github.com/vmware-labs/yaml-jsonpointer"
	yaml "gopkg.in/yaml.v3"
)

func ChangeYamlNodeValue(keyPath string, setValue string) {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	var yamlNode yaml.Node
	yaml.Unmarshal([]byte(yamlFile), &yamlNode)

	r, _ := yptr.Find(&yamlNode, keyPath)
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)
	r.SetString(setValue)
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)

	f, err := os.Create("values.yaml")
	if err != nil {
		log.Fatalf("Problem opening file: %v", err)
		log.Fatalf("Problem opening file: %v", err)
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	encoder.Encode(yamlNode.Content[0])
	encoder.Close()
}
