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

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating a new user...")
	type createUserParams struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := createUserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing request body: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:   uuid.New(),
		Name: params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Creating user failed with error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUserModel(user))
	log.Printf("New user %s created successfully.", user.Name)
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	log.Printf("Getting a user...")
	respondWithJSON(w, http.StatusOK, databaseUserToUserModel(user))
	log.Printf("User retrieved successfully with the given API Key")
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	log.Printf("Getting all posts for the user %s...", user.Name)
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  5,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Getting posts for user failed with error: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPostsModel(posts))
	log.Printf("Posts for user %s retrieved successfully.", user.Name)
}