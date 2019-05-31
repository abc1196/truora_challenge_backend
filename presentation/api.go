package presentation

import (
	"log"
	"net/http"

	"../utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// handleRequests - Setups the API Rest with the two needed endpoints
func handleRequests() {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Route("/api/challenge/v1/domains", func(r chi.Router) {
		r.Get("/{id}", CreateDomain)
		r.Get("/", GetDomains)
	})
	log.Fatal(http.ListenAndServe(":4000", r))
}

// Init - Main method to connect to the database and setup the API Rest
func Init() {
	utils.Connect()
	handleRequests()
}
