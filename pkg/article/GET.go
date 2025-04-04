package article

import (
	"StyleSwap/pkg/filter"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (config *ArticleConfig) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire les filtres
	filtre := filter.ParseFilterCriteria(r)
	articles, err := config.ArticleRepository.FindByCriteria(filtre)
	utils.HandleError(err, "Error while fetching articles from database")

	// Les données de réponse
	var res []model.ArticleResponse
	for _, article := range articles {
		var imageBase64 string
		if article.ImageData != nil { // ✅ Évite d'encoder nil en base64
			imageBase64 = base64.StdEncoding.EncodeToString(article.ImageData)
		}

		res = append(res, model.ArticleResponse{
			ArticleId:          int(article.ID),
			UserPseudo:         article.PseudoUser,
			ArticleName:        article.Name,
			ArticlePrice:       article.Price,
			ArticleSize:        article.Size,
			ArticleBrand:       article.Brand,
			ArticleColor:       article.Color,
			ArticleDescription: article.Description,
			ArticleImage:       imageBase64, // ✅ Gère le cas où l'image est absente
		})
	}

	render.JSON(w, r, res)
}

// GetArticleID gère la récupération d'un article en fonction de l'ID de l'article

func (config *ArticleConfig) GetArticleID(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'article à partir des paramètres de l'URL
	ArticleID := chi.URLParam(r, "id")
	if ArticleID == "" {
		render.JSON(w, r, map[string]string{"error": "Article ID is required"})
		return
	}
	ArticleIDInt, err := strconv.Atoi(ArticleID)
	utils.HandleError(err, "Error while converting article ID to integer")

	// Rechercher l'article dans la base de données
	article, err := config.ArticleRepository.FindByID(ArticleIDInt)
	utils.HandleError(err, "Error while fetching article from database")

	// Vérifier que le User est autorisé à voir l'article
	VerifArticle(config, article, w, r)

	//Les données de réponse
	res := model.ArticleResponse{
		ArticleId:          int(article.ID),
		UserPseudo:         article.PseudoUser,
		ArticleName:        article.Name,
		ArticlePrice:       article.Price,
		ArticleSize:        article.Size,
		ArticleBrand:       article.Brand,
		ArticleDescription: article.Description,
		ArticleImage:       string(article.ImageData),
	}
	render.JSON(w, r, res)
}
