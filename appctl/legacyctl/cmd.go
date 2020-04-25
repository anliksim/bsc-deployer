package legacyctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/anliksim/bsc-deployer/appctl"
	"github.com/anliksim/bsc-deployer/appctl/kubectl"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"log"
	"net/http"
	"strings"
)

const appPath = "/home/anliksim/codebase/bsc-env/apps"

func Apply() {
	jsonString := kubectl.GetLegacyDescriptorsAsJson(appPath)
	// multiple Deployments are returned as part of kind List by kubectl
	if strings.Contains(jsonString, "List") {
		ForEachDeploymentInList([]byte(jsonString), func(payload []byte) {
			runDeployment(payload)
		})
	} else {
		runDeployment([]byte(jsonString))
	}
}

func Delete() {
	jsonString := kubectl.GetLegacyDescriptorsAsJson(appPath)
	// multiple Deployments are returned as part of kind List by kubectl
	if strings.Contains(jsonString, "List") {
		ForEachDeploymentInList([]byte(jsonString), func(payload []byte) {
			runStop(payload)
		})
	} else {
		runStop([]byte(jsonString))
	}
}

func ForEachDeploymentInList(jsonContent []byte, deploymentHandler func([]byte)) {
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

func runStop(payload []byte) {
	deployment := JsonToDeployment(payload)
	name := deployment.Name
	host := deployment.Spec.Template.Annotations["legacy/host"]
	log.Printf("Deleting apps from %s...", host)
	deleteProcess(host, name)
}

func runDeployment(payload []byte) {
	deployment := JsonToDeployment(payload)
	host := deployment.Spec.Template.Annotations["legacy/host"]
	log.Printf("Deploying to %s...", host)
	postProcesses(host, payload)
}

func postProcesses(host string, payload []byte) {
	call(func() (response *http.Response, e error) {
		return http.Post(serverUrl(host, "processes"), "application/json", bytes.NewBuffer(payload))
	}, func(body []byte) {
		printResponse(body)
	})
}

func deleteProcess(host string, name string) {
	call(func() (response *http.Response, e error) {
		req, err := http.NewRequest(http.MethodDelete, serverUrl(host, fmt.Sprintf("processes/%s", name)), nil)
		if err != nil {
			log.Fatalf("Error on delete: %v", err)
		}
		return http.DefaultClient.Do(req)
	}, func(body []byte) {
		printResponse(body)
	})
}

func printResponse(body []byte) {
	appctl.SetDarkGray()
	fmt.Printf("%s\n", body)
	appctl.SetNoColor()
}

func call(httpCall func() (*http.Response, error), callback func([]byte)) {
	resp, err := httpCall()
	if err != nil {
		log.Fatalf("Error on get: %v", err)
	}
	handle(resp, callback)
}

func handle(resp *http.Response, callback func([]byte)) {
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error parsing response: %v", err)
		}
		callback(body)
	}
}

func serverUrl(host string, path string) string {
	return fmt.Sprintf("%s/%s", host, path)
}
