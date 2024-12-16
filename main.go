package main

import (
	"StyleSwap/config"
	"StyleSwap/pkg/article"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Mount("/api", article.Routes(configuration))
	return router
}

func main() {
	configuration, err := config.New()
	if err != nil {
		log.Panicln("Configuration error:",err)
	}

	router := Routes(configuration)
	
	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}	