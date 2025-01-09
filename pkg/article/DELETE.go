package article

import (
	"StyleSwap/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (config *ArticleConfig) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	idArticle := chi.URLParam(r, "id4Delete")
	if idArticle == "" {
		json.NewEncoder(w).Encode("Article ID is required")
		return
	}
	id, err := strconv.Atoi(idArticle)
	utils.HandleError(err, "Error while converting article ID to integer")

	article, err := config.ArticleRepository.FindByID(id)
	utils.HandleError(err, "Error while fetching article from database")

	if article != nil {
		fmt.Println("Article found")
		errBucket := utils.DeleteImageFromS3(article.ImageURL)
		utils.HandleError(errBucket, "Error while deleting image from S3")
	} else {
		json.NewEncoder(w).Encode("Article not found")
		return // Return if article not found
	}

	errBDD := config.ArticleRepository.Delete(int(article.ID))
	utils.HandleError(errBDD, "Error while deleting article from database")
	if errBDD == nil {
		json.NewEncoder(w).Encode("Article deleted successfully")
	}
}
