package apiv1

import (
	"github.com/anliksim/bsc-deployer/api"
	"github.com/anliksim/bsc-deployer/appctl/kubectl"
	"github.com/anliksim/bsc-deployer/appctl/legacyctl"
	"github.com/anliksim/bsc-deployer/config"
	"github.com/anliksim/bsc-deployer/model"
	modelv1 "github.com/anliksim/bsc-deployer/model/v1"
	"github.com/anliksim/bsc-deployer/util"
	"github.com/gorilla/mux"
	"github.com/nvellon/hal"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var baseUrl string

var deployments = make(map[time.Time]string)

func Register(r *mux.Router, base string) {
	baseUrl = base
	r.HandleFunc(Path(""), getBase)
	r.HandleFunc(Path(api.Health), getHealth)
	r.HandleFunc(Path(api.Deployments), getDeploy).Methods("GET")
	r.HandleFunc(Path(api.Deployments), postDeploy).Methods("POST")
	r.HandleFunc(Path(api.Deployments), deleteDeploy).Methods("DELETE")
}

func getBase(w http.ResponseWriter, r *http.Request) {
	res := hal.NewResource(&model.None{}, Url(baseUrl, ""))
	res.AddNewLink("health", Url(baseUrl, api.Health))
	res.AddNewLink("deployments", Url(baseUrl, api.Deployments))
	util.RespondJson(w, res)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	res := hal.NewResource(&modelv1.Health{
		Status: "ok",
	}, Url(baseUrl, api.Health))
	util.RespondJson(w, res)
}

func getDeploy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requesting deployment status")
	res := hal.NewResource(&modelv1.Deployments{
		Entries: deployments,
	}, Url(baseUrl, api.Deployments))
	util.RespondJson(w, res)
}

func postDeploy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request for new deployment")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}
	deployData := config.ParseJson(body)
	log.Printf("%v\n", deployData)

	now := time.Now()
	deployments[now] = "apply: " + deployData.Rev

	// async
	go deploy()

	util.Respond(w, now.Format("2006-01-02 15:04:05"))
}

func deleteDeploy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request for new deployment")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}
	deployData := config.ParseJson(body)
	log.Printf("%v\n", deployData)

	now := time.Now()
	deployments[now] = "delete: " + deployData.Rev

	// async
	go deleteAll()

	util.Respond(w, now.Format("2006-01-02 15:04:05"))
}

func deploy() {
	kubectl.SetContext("minikube")
	log.Print("Version:")
	kubectl.ShortVersion()
	kubectl.DeployPolicies()
	kubectl.DeployAppsForPrivate()
	legacyctl.Apply()
}

func deleteAll() {
	legacyctl.Delete()
}
