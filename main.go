package main

import (
	"fmt"
	"log"
	"os"

	// "net/http"
	// git "github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing/transport/http"
	// "github.com/xanzy/go-gitlab"
	// "gopkg.in/yaml.v3"
	gin "github.com/gin-gonic/gin"
	yptr "github.com/vmware-labs/yaml-jsonpointer"
	yaml "gopkg.in/yaml.v3"
)

func main() {

	yamlNode := changeYamlNodeValue("/uhc/tag", "12314")

	f, err := os.Create("values.yaml")
	if err != nil {
		log.Fatalf("Problem opening file: %v", err)
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	encoder.Encode(yamlNode.Content[0])
	encoder.Close()

	router := gin.Default()
	router.GET("/values", getValues)

	router.Run(":8080")

}

func changeYamlNodeValue(keyPath string, setValue string) (yamlNode yaml.Node) {

	yamlFile, err := os.ReadFile("values.yaml")
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal([]byte(yamlFile), &yamlNode)

	r, _ := yptr.Find(&yamlNode, keyPath)
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)
	r.SetString(setValue)
	fmt.Printf("Scalar %q at %d:%d\n", r.Value, r.Line, r.Column)

	return yamlNode
}

func getValues(c *gin.Context) {
	c.File("values.yaml")
}
