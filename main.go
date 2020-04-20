package main

import (
	"github.com/anliksim/bsc-deployer/config"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var deployments = make(map[time.Time]string)

func main() {
	errorChain := alice.New(loggerHandler, recoverHandler)
	r := mux.NewRouter()
	http.Handle("/", errorChain.Then(r))
	r.HandleFunc("/", getMain)
	r.HandleFunc("/health", getHealth)
	r.HandleFunc("/deployments", getDeploy).Methods("GET")
	r.HandleFunc("/deployments", postDeploy).Methods("POST")
	log.Printf("Starting server at localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Error starting legacyctld: %v", err)
	}
}

func getMain(w http.ResponseWriter, r *http.Request) {
	respond(w, "deployer")
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	respond(w, "ok")
}

func getDeploy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Requesting deployment status")
	respond(w, deploymentsAsString())
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
	deployments[now] = deployData.Rev

	respond(w, now.String())
}

func loggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(">> %s %s", r.Method, r.URL.Path)
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func deploymentsAsString() string {
	var str strings.Builder
	dl := " scheduled at "
	for key, value := range deployments {
		str.WriteString(value)
		str.WriteString(dl)
		str.WriteString(key.String())
		str.WriteRune('\n')
	}
	return str.String()
}

func respond(w http.ResponseWriter, message string) {
	if _, err := io.WriteString(w, message); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
