package templatedata

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

func (t *TemplateData) EncodeGitlabTemplate() (string, error) {

	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	defer encoder.Close()

	encoder.SetIndent(2)
	err := encoder.Encode(t.GitlabTemplate)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (t *TemplateData) EncodeChart() (string, error) {

	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	defer encoder.Close()

	encoder.SetIndent(2)
	err := encoder.Encode(t.Chart)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (t *TemplateData) EncodeValues() (string, error) {

	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	defer encoder.Close()

	encoder.SetIndent(2)
	err := encoder.Encode(t.Values)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (t *TemplateData) PopulateArgotTemplate() (string, error) {

	serviceName := t.Chart.Name
	k8sRepo := t.K8sRepo
	tmpl, err := os.ReadFile("templates/argocdtTemplate.yaml")

	if err != nil {
		return "", err
	}

	tmpl = bytes.ReplaceAll(tmpl, []byte("replacemetemplate"), []byte(serviceName))
	tmpl = bytes.ReplaceAll(tmpl, []byte("replacemeproject"), []byte(k8sRepo))
	var data bytes.Buffer
	data.Write(tmpl)
	at := data.String()
	return at, nil

}
