package article

import (
	"StyleSwap/database/dbmodel"
	"StyleSwap/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// UpdateArticleHandler gère la mise à jour d'un article en fonction de l'ID de l'article

func (config *ArticleConfig) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	idstring := chi.URLParam(r, "id4Update")
	id, errConv := strconv.Atoi(idstring)
	utils.HandleError(errConv, "Error while converting article ID to integer")

	article_name := r.FormValue("name")
	article_price_str := r.FormValue("price")
	article_price, errConv := strconv.Atoi(article_price_str)
	utils.HandleError(errConv, "Error while converting price to integer")
	article_size := r.FormValue("size")
	article_brand := r.FormValue("brand")
	article_color := r.FormValue("color")
	article_description := r.FormValue("description")
	_, _, err := r.FormFile("image")

	imageURL := ""

	if err == nil {
		// Upload the image to S3
		file, _, err := r.FormFile("image")
		utils.HandleError(err, "Error retrieving image from form data")
		defer file.Close()

		// Générer un nom unique pour le fichier sur S3
		uniqueID := uuid.New().String()
		filename := fmt.Sprintf("%s.png", uniqueID) // Utiliser un nom unique pour éviter les collisions

		// Télécharger l'image sur S3
		imageURL, err = utils.UploadToS3(file, filename)
		utils.HandleError(err, "Error uploading image to S3")
	}

	article := &dbmodel.ArticleEntry{
		Name:        article_name,
		Price:       article_price,
		Size:        article_size,
		Brand:       article_brand,
		Color:       article_color,
		Description: article_description,
		ImageURL:    imageURL,
	}
	errUpdate := config.ArticleRepository.Update(article, id)
	utils.HandleError(errUpdate, "Error while updating article in database")

	render.JSON(w, r, map[string]string{"message": "Article updated successfully"})
}
