package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/kinshuk1993/rssaggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// loads the .env by default, but can specify file name as param too
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not set")
	}

	dbUrlString := os.Getenv("DB_URL")
	if dbUrlString == "" {
		log.Fatal("DB_URL is not set")
	}

	conn, err := sql.Open("postgres", dbUrlString)
	if err != nil {
		log.Fatal("Cannot connect to the database: ", err)
	}

	apiCfg := apiConfig {
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetAllFeeds)

	router.Mount("/v1", v1Router)


	srv := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf(":%s", portString),
	}

	log.Printf("Server starting on %s", portString)
	errStartingServer := srv.ListenAndServe()
	if errStartingServer != nil  {
		log.Fatal(errStartingServer)
	}

	fmt.Println("PORT is set to:", portString)
}
