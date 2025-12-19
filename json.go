package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with error: ", message)
	}

	// this would be something like: {"error": "message here"} - basically specifying that the key in the marshalled json is "error" and similarly, unmarshall for error string takes the value of the key error and puts it in the Error field of the struct
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// data returned is in a byte format so that we can return it easily into a http response
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}