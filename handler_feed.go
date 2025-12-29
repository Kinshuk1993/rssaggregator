package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kinshuk1993/rssaggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	log.Printf("Creating a new feed...")
	type createFeed struct {
		Name string `json:"name"`
		URL string `json:"URL"`
	}
	decoder := json.NewDecoder(r.Body)
	params := createFeed{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:   uuid.New(),
		Name: params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: params.URL,
		Userid: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Creating user failed with error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeedModel(feed))
	log.Printf("User %s created a new feed with ID %s successfully.", user.Name, feed.ID)
}
