package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ArticleConfig struct {
	*config.Config
}

func New(configuration *config.Config) *ArticleConfig {
	return &ArticleConfig{configuration}
}

// ArticleHandler gère la création d'un nouvel article

func (config *ArticleConfig) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse la requête multipart pour obtenir les données de formulaire
	err := r.ParseMultipartForm(10 << 20) // Limite de taille de 10 Mo pour l'image
	utils.HandleError(err, "Error parsing form data")
	// Récupérer le pseudo de l'utilisateur
	userpseudo := r.FormValue("userPseudo")
	if userpseudo == "" {
		render.JSON(w, r, map[string]string{"error": "Missing user pseudo"})
		return
	}
	user, err := config.UserRepository.FindByPseudo(userpseudo)
	utils.HandleError(err, "Error while fetching user from the database")

	// Extraire les données de l'article du formulaire
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	utils.HandleError(err, "Error while converting price to integer")
	size := r.FormValue("size")
	brand := r.FormValue("brand")
	color := r.FormValue("color")
	description := r.FormValue("description")

	// Vérification que tous les champs sont fournis
	if name == "" || price == 0 || description == "" {
		render.JSON(w, r, map[string]string{"error": "Missing required fields"})
		return
	}

	// Récupérer le fichier de la requête (champ "image")
	file, _, err := r.FormFile("image")
	utils.HandleError(err, "Error retrieving image from form data")
	defer file.Close()

	// Générer un nom unique pour le fichier sur S3
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s.png", uniqueID) // Utiliser un nom unique pour éviter les collisions

	// Télécharger l'image sur S3
	imageURL, err := utils.UploadToS3(file, filename)
	utils.HandleError(err, "Error uploading image to S3")

	// Créer un nouvel article dans la base de données
	articleEntry := &dbmodel.ArticleEntry{
		PseudoUser:  user.Pseudo,
		Name:        name,
		Price:       price,
		Size:        size,
		Brand:       brand,
		Color:		 color,
		Description: description,
		ImageURL:    imageURL,
	}

	// Ajouter l'article à la base de données
	if err := config.ArticleRepository.Create(articleEntry); err != nil {
		render.JSON(w, r, map[string]string{"error": "Error while adding article to the database"})
		return
	}

	// Répondre avec un message de succès
	render.JSON(w, r, map[string]string{"message": "Article added successfully", "image_url": imageURL})
}



// GetArticlesHandler gère la récupération de tous les articles ou d'un article spécifique en fonction du paramètre UserPseudo


func (config *ArticleConfig) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	Pseudo := r.URL.Query().Get("UserPseudo")
	var articles []dbmodel.ArticleEntry
	var err error

	// Recherche des articles en fonction du paramètre UserPseudo
	if Pseudo != "" {
		articles, err = config.ArticleRepository.FindByPseudo(Pseudo)
	} else {
		articles, err = config.ArticleRepository.FindAll()
	}

	// Gestion des erreurs si la recherche échoue
	utils.HandleError(err, "Error while fetching articles from the database")

	// Créer un tableau pour contenir toutes les réponses d'articles
	var articleResponses []model.ArticleResponse

	// Préparer les données de réponse
	for _, article := range articles {
		res := model.ArticleResponse{
			UserPseudo:         article.PseudoUser,
			ArticleName:        article.Name,
			ArticlePrice:       article.Price,
			ArticleSize:        article.Size,
			ArticleBrand:       article.Brand,
			ArticleDescription: article.Description,
			ArticleImage:       article.ImageURL,
		}
		articleResponses = append(articleResponses, res)
	}

	// Renvoyer tous les articles dans une seule réponse
	render.JSON(w, r, articleResponses)
}







// DeleteArticleHandler gère la suppression d'un article en fonction de l'ID de l'article

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
	_, _,err := r.FormFile("image")

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




// GetArticleID gère la récupération d'un article en fonction de l'ID de l'article

func (config *ArticleConfig) GetArticleID(w http.ResponseWriter, r *http.Request){
	// Récupérer l'ID de l'article à partir des paramètres de l'URL
	ArticleID := chi.URLParam(r, "id")
	ArticleIDInt, err := strconv.Atoi(ArticleID)
	utils.HandleError(err, "Error while converting article ID to integer")

	// Rechercher l'article dans la base de données
	article, err := config.ArticleRepository.FindByID(ArticleIDInt)
	utils.HandleError(err, "Error while fetching article from database")

	//Les données de réponse
	res := model.ArticleResponse{
		ArticleId:          int(article.ID),
		UserPseudo:         article.PseudoUser,
		ArticleName:        article.Name,
		ArticlePrice:       article.Price,
		ArticleSize:        article.Size,
		ArticleBrand:       article.Brand,
		ArticleDescription: article.Description,
		ArticleImage:       article.ImageURL,
	}
	render.JSON(w, r, res)
}
