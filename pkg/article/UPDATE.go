package article

import (
	"StyleSwap/database/dbmodel"
	"StyleSwap/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// UpdateArticleHandler gère la mise à jour d'un article en fonction de l'ID de l'article

func (config *ArticleConfig) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	idstring := chi.URLParam(r, "id4Update")
	if idstring == "" {
		render.JSON(w, r, map[string]string{"error": "Article ID is required"})
		return
	}
	id, errConv := strconv.Atoi(idstring)
	utils.HandleError(errConv, "Error while converting article ID to integer")

	article, errArticle := config.ArticleRepository.FindByID(id)
	utils.HandleError(errArticle, "Error while fetching article from database")

	// Vérifier que l'utilisateur est autorisé à modifier l'article
	VerifArticle(config, article, w, r)

	article_name := r.FormValue("name")
	article_price_str := r.FormValue("price")

	var article_price int
	if article_price_str != "" {
		article_price, errConv = strconv.Atoi(article_price_str)
		utils.HandleError(errConv, "Error while converting price to integer")
	}
	article_size := r.FormValue("size")
	article_brand := r.FormValue("brand")
	article_color := r.FormValue("color")
	article_description := r.FormValue("description")

	article = &dbmodel.ArticleEntry{
		Name:        article_name,
		Price:       article_price,
		Size:        article_size,
		Brand:       article_brand,
		Color:       article_color,
		Description: article_description,
	}
	errUpdate := config.ArticleRepository.Update(article, id)
	utils.HandleError(errUpdate, "Error while updating article in database")

	render.JSON(w, r, map[string]string{"message": "Article updated successfully"})
}
