package main

import (
	"StyleSwap/config"
	"StyleSwap/pkg/article"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/chat"
	"StyleSwap/pkg/usermanagement"
	"StyleSwap/pkg/usermanagement/login"
	"StyleSwap/pkg/usermanagement/register"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(configuration *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// Middleware CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes publiques
	router.Mount("/api/v1/register", register.Routes(configuration))
	router.Mount("/api/v1/login", login.Routes(configuration))

	// Routes protégées
	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware("StyleSwap"))
		r.Mount("/api/v1/articles", article.Routes(configuration))
		r.Mount("/api/v1/chat", chat.Routes(configuration))
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
