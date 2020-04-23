package util

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Respond(w http.ResponseWriter, message string) {
	if _, err := io.WriteString(w, message); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func RespondJson(w http.ResponseWriter, v interface{}) {
	payload, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}
	if _, err := io.WriteString(w, string(payload)); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
