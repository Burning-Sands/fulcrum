package main

import (
	"fmt"
	"log"
	"os"

	"net/http"
	// git "github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing/transport/http"
	// "github.com/xanzy/go-gitlab"
	// "gopkg.in/yaml.v3"
	gin "github.com/gin-gonic/gin"
	yptr "github.com/vmware-labs/yaml-jsonpointer"
	yaml "gopkg.in/yaml.v3"
)

func main() {

	router := gin.Default()
	router.Static("/", "./public")
	//router.GET("/values", getValues)
	router.POST("/upload", postValues)

	router.Run(":8080")

}

func changeYamlNodeValue(keyPath string, setValue string) {
	
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
	}
	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	encoder.Encode(yamlNode.Content[0])
	encoder.Close()
}

func getValues(c *gin.Context) {
	c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
}

func postValues(c *gin.Context) {

	path := c.PostForm("path")
	value := c.PostForm("value")
	changeYamlNodeValue(path, value)

	// Source
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	c.String(http.StatusBadRequest, "get form err: %s", err.Error())
	// 	return
	// }

	// if err := c.SaveUploadedFile(file, "components.yaml"); err != nil {
	// 	c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
	// 	return
	// }

	c.String(http.StatusOK, "File values.yaml updated successfully with path=%s and value=%s.", path, value )
}
