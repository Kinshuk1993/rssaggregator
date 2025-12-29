package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kinshuk1993/rssaggregator/internal/auth"
	"github.com/kinshuk1993/rssaggregator/internal/database"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, dbUser database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(&r.Header)

		if err != nil {
			respondWithError(w, http.StatusForbidden, err.Error())
			log.Printf("Error getting API Key from request header: %v", err)
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("no user found using the given api key: %v", err))
			log.Printf("no user found using the given api key: %v", err)
			return
		}

		handler(w, r, user)
	}
}