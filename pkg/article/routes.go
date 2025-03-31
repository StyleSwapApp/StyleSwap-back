package article

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux {
	articleConfig := New(configuration)
	router := chi.NewRouter()

	router.Post("/newArticle", articleConfig.ArticleHandler)
	router.Get("/getArticles", articleConfig.GetArticlesHandler)
	router.Delete("/{id4Delete}", articleConfig.DeleteArticleHandler)
	router.Put("/{id4Update}", articleConfig.UpdateArticleHandler)
	router.Get("/{id}", articleConfig.GetArticleID)

	return router
}
