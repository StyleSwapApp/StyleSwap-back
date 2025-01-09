package article

import (
	"StyleSwap/config"

	"github.com/go-chi/chi/v5"
)

func Routes(configuration *config.Config) *chi.Mux{
	articleConfig := New(configuration)
	router := chi.NewRouter()

	router.Post("/newArticle", articleConfig.ArticleHandler)
	router.Get("/getArticle", articleConfig.GetArticlesHandler)
	router.Delete("/deleteArticle", articleConfig.DeleteArticleHandler)
	router.Put("/updateArticle", articleConfig.UpdateArticleHandler)
	router.Get("/{id}", articleConfig.GetArticleID)
	

	return router
}