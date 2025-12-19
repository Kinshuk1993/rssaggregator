package main

import (
	"log"
	"net/http"
)

func handlerError(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Checking error status...")
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	log.Printf("Error check completed.")
}