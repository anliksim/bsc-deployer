package test

import (
	"bytes"
	"fmt"
	"github.com/anliksim/bsc-deployer/config"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

const protocol = "http"
const host = "localhost"
const port = "3557"

const request = `
{
	"dir": "/home/anliksim/codebase/bsc-env",
	"rev": "ff755b0 Implement policies"
}
`

func TestDeployments(t *testing.T) {
	deploymentRequest := config.ParseJson([]byte(request))
	log.Printf("%v\n", deploymentRequest)

	postDeployments([]byte(request))
	getDeployments()
}

func TestDeleteLegacy(t *testing.T) {
	deploymentRequest := config.ParseJson([]byte(request))
	log.Printf("%v\n", deploymentRequest)

	deleteDeployments([]byte(request))
	getDeployments()
}

func getDeployments() {
	call(func() (response *http.Response, e error) {
		return http.Get(deploymentsUrl())
	}, func(body []byte) {
		fmt.Printf("%s\n", body)
	})
}

func postDeployments(payload []byte) {
	call(func() (response *http.Response, e error) {
		return http.Post(deploymentsUrl(), "application/json", bytes.NewBuffer(payload))
	}, func(body []byte) {
		fmt.Printf("%s\n", body)
	})
}

func deleteDeployments(payload []byte) {
	call(func() (response *http.Response, e error) {
		req, err := http.NewRequest(http.MethodDelete, deploymentsUrl(), bytes.NewBuffer(payload))
		if err != nil {
			log.Fatalf("Error on delete: %v", err)
		}
		return http.DefaultClient.Do(req)
	}, func(body []byte) {
		fmt.Printf("%s\n", body)
	})
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

func deploymentsUrl() string {
	return serverUrl("v1/deployments")
}

func serverUrl(path string) string {
	return fmt.Sprintf("%s://%s:%s/%s", protocol, host, port, path)
}
