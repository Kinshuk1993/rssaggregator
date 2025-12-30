package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kinshuk1993/rssaggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollowsByUserID(w http.ResponseWriter, r *http.Request, user database.User) {
	log.Printf("Getting all feeds followed by user %s...", user.Name)

	allFeedsByUser, err := apiCfg.DB.GetFeedFollowsByUserID(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Getting all feeds for user %s failed with error: %v", user.Name, err))
		return
	}

	// if len(allFeedsByUser) == 0 {
	// 	respondWithError(w, http.StatusNotFound, fmt.Sprintf("No feeds found followed by user %s", user.Name))
	// 	return
	// }

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowsByUserIDToFeedFollowsByUserIDModel(allFeedsByUser))
	log.Printf("User %s follows %d feeds.", user.Name, len(allFeedsByUser))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollowsByUserID(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDString := chi.URLParam(r, "feedFollowID")
	feedFollowDeleteID, err := uuid.Parse(feedFollowIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Parsing feedFollowID %s failed with error: %v", feedFollowIDString, err))
		return
	}

	log.Printf("Deleting a feed with id %v followed by user %s...", feedFollowDeleteID, user.Name)

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowDeleteID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error deleting feed for user %s failed with error: %v", user.Name, err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
	log.Printf("User %s deleted the feed with id %s.", user.Name, feedFollowIDString)
}