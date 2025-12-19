package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Checking readiness status...")
	respondWithJSON(w, http.StatusOK, struct{}{})
	log.Printf("Readiness check passed.")
}
