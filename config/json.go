package config

import (
	"encoding/json"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"log"
)

type DeploymentData struct {
	Dir string `json:"dir"`
	Rev string `json:"rev"`
}

type NamedObject struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
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

func ForEachItemInList(jsonContent []byte, deploymentHandler func([]byte)) {
	list := new(v1.List)
	if err := json.Unmarshal(jsonContent, &list); err != nil {
		log.Fatalf("Error during unmarshal: %v", err)
	}
	for _, c := range list.Items {
		deploymentHandler(c.Raw)
	}
}

func JsonToDeployment(jsonContent []byte) *appsv1.Deployment {
	deployment := new(appsv1.Deployment)
	if err := json.Unmarshal(jsonContent, &deployment); err != nil {
		log.Fatalf("Error during unmarshal: %v", err)
	}
	return deployment
}

func JsonToNamedObject(jsonContent []byte) *NamedObject {
	named := new(NamedObject)
	if err := json.Unmarshal(jsonContent, &named); err != nil {
		log.Fatalf("Error during unmarshal: %v", err)
	}
	return named
}
