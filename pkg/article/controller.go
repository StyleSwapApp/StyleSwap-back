package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
	"net/http"

	"github.com/go-chi/render"
)

type ArticleConfig struct {
	*config.Config
}

func New(configuration *config.Config) *ArticleConfig {
	return &ArticleConfig{configuration}
}

func (conifg *ArticleConfig) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.ArticleRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}

	articleEntry := &dbmodel.ArticleEntry{Name: req.Name, Price: req.Price, Description: req.Description}
	if err := conifg.ArticleRepository.Create(articleEntry); err != nil {
		render.JSON(w, r, map[string]string{"error": "Error while adding article to the database"})
		return
	}

	render.JSON(w, r, "Bien ajouté à la BDD")
}

func (config *ArticleConfig) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := config.ArticleRepository.FindAll()
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Error while fetching articles from the database"})
		return
	}

	render.JSON(w, r, articles)
}
