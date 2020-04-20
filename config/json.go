package config

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
)

type DeploymentData struct {
	Dir string `json:"dir"`
	Rev string `json:"rev"`
}

func UnmarshalFile(filePath string) *DeploymentData {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return ParseJson(content)
}

func ParseJson(jsonContent []byte) *DeploymentData {
	deployment := new(DeploymentData)
	if err := json.Unmarshal(jsonContent, &deployment); err != nil {
		log.Fatal(err)
	}
	return deployment
}
