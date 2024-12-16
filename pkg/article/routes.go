package article

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux{
	articleConfig := New(configuration)
	router := chi.NewRouter()

	router.Get("/articles", articleConfig.ArticleHandler)

	return router
}