package yamleditor

import (
	"bytes"
	"fmt"

	yptr "github.com/vmware-labs/yaml-jsonpointer"
	yaml "gopkg.in/yaml.v3"
)

type YamlOperator struct {
	YamlNode yaml.Node
	Buffer   bytes.Buffer
}

func ChangeYamlNodeValue(node *yaml.Node, keyPath string, setValue string) {

	r, err := yptr.Find(node, keyPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found scalar %q at %d:%d\n", r.Value, r.Line, r.Column)
	r.SetString(setValue)
	fmt.Printf("Set scalar to %q at %d:%d\n", r.Value, r.Line, r.Column)

}

func IterateOverYamlNode(node *yaml.Node, buffer *bytes.Buffer) {
	for _, v := range node.Content {
		switch v.Kind {
		case yaml.SequenceNode:
			IterateOverYamlNode(v, buffer)
		case yaml.MappingNode:
			IterateOverYamlNode(v, buffer)
		case yaml.ScalarNode:
			//fmt.Println(v.Value)
			buffer.WriteString(string(v.Value) + " ")
		default:
			fmt.Println("Unknown node kind:", node.Kind)
		}
	}
}
