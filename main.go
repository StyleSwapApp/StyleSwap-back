package main

import (
	"StyleSwap/config"
	"StyleSwap/pkg/article"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/usermanagement"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// Appliquer les middleware globalement ou sur un groupe de routes non protégées ici
	router.Mount("/api/v1/login", usermanagement.Routes(configuration)) // Route non protégée

	// Protéger les routes spécifiques avec le middleware d'authentification
	router.Group(func(r chi.Router) { 
		r.Use(auth.AuthMiddleware("StyleSwap"))
		r.Mount("/api/v1/articles", article.Routes(configuration))
		r.Mount("/api/v1/user", usermanagement.Routes(configuration))
	})

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
