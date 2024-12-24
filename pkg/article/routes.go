package article

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux{
	articleConfig := New(configuration)
	router := chi.NewRouter()

	router.Post("/newArticle", articleConfig.ArticleHandler)

	return router
}