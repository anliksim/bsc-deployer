package main

import (
	"github.com/anliksim/bsc-deployer/api"
	apiv1 "github.com/anliksim/bsc-deployer/api/v1"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

const version = "v0.8"
const port = ":3557"
const baseUrl = "http://localhost" + port

func main() {
	errorChain := alice.New(loggerHandler, recoverHandler)
	r := mux.NewRouter()
	http.Handle(api.Base, errorChain.Then(r))
	api.Register(r, baseUrl)
	apiv1.Register(r, baseUrl)
	log.Printf("Starting server %s at %s", version, baseUrl)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting deployer: %v", err)
	}
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
