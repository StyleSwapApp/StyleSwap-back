package article

import (
	"StyleSwap/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (config *ArticleConfig) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	idArticle := chi.URLParam(r, "id4Delete")
	id, err := strconv.Atoi(idArticle)
	utils.HandleError(err, "Error while converting article ID to integer")

	article, err := config.ArticleRepository.FindByID(id)
	utils.HandleError(err, "Error while fetching article from database")

	errBucket := utils.DeleteImageFromS3(article.ImageURL)
	utils.HandleError(errBucket, "Error while deleting image from S3")

	errBDD := config.ArticleRepository.Delete(int(article.ID))
	utils.HandleError(errBDD, "Error while deleting article from database")
}
