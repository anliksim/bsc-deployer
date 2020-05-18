package legacyctl

import (
	"bytes"
	"fmt"
	"github.com/anliksim/bsc-deployer/appctl/kubectl"
	"github.com/anliksim/bsc-deployer/config"
	"github.com/anliksim/bsc-deployer/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func appsPath(dirPath string) string {
	return dirPath + "/apps"
}

func Apply(dirPath string) {
	jsonString := kubectl.GetLegacyDescriptorsAsJson(appsPath(dirPath))
	// multiple Deployments are returned as part of kind List by kubectl
	if strings.Contains(jsonString, "List") {
		config.ForEachItemInList([]byte(jsonString), func(payload []byte) {
			runDeployment(payload)
		})
	} else {
		runDeployment([]byte(jsonString))
	}
}

func Delete(dirPath string) {
	jsonString := kubectl.GetLegacyDescriptorsAsJson(appsPath(dirPath))
	// multiple Deployments are returned as part of kind List by kubectl
	if strings.Contains(jsonString, "List") {
		config.ForEachItemInList([]byte(jsonString), func(payload []byte) {
			runStop(payload)
		})
	} else {
		runStop([]byte(jsonString))
	}
}

func runStop(payload []byte) {
	deployment := config.JsonToDeployment(payload)
	name := deployment.Name
	host := deployment.Spec.Template.Annotations["legacy/host"]
	log.Printf("Deleting apps from %s...", host)
	deleteProcess(host, name)
}

func runDeployment(payload []byte) {
	deployment := config.JsonToDeployment(payload)
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
	util.SetDarkGray()
	fmt.Printf("%s\n", body)
	util.SetNoColor()
}

func call(httpCall func() (*http.Response, error), callback func([]byte)) {
	resp, err := httpCall()
	if err != nil {
		log.Printf("Error on get: %v", err)
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
