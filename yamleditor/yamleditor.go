package yamleditor

import (
	"os"

	yqlib "github.com/mikefarah/yq/v4/pkg/yqlib"
	yaml "gopkg.in/yaml.v3"
)

func ChangeYamlNodeValue(keyPath string, setValue string) {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	m := make(map[string]*yqlib.CandidateNode)
	var yamlNode yaml.Node
	var candidateNode yqlib.CandidateNode
	yaml.Unmarshal([]byte(yamlFile), &yamlNode)
	candidateNode.UnmarshalYAML(&yamlNode, m)

	yqlib.NodeToString(&candidateNode)

	// r, _ := yptr.Find(&yamlNode, keyPath)
	// fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)
	// r.SetString(setValue)
	// fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)

	// f, err := os.Create("values.yaml")
	// if err != nil {
	// 	log.Fatalf("Problem opening file: %v", err)
	// }
	// encoder := yqlib.NewYamlEncoder(2, false)
	// encoder.Encode(f, yamlNode.Content[0])

}
