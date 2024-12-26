package main

import (
	"StyleSwap/config"
	"StyleSwap/pkg/article"
	"StyleSwap/pkg/usermanagement"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Mount("/api/v1/articles", article.Routes(configuration))
	router.Mount("/api/v1/user", usermanagement.Routes(configuration))
	return router
}

func main() {
	configuration, err := config.New()
	if err != nil {
		log.Panicln("Configuration error:", err)
	}

	router := Routes(configuration)

	log.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
