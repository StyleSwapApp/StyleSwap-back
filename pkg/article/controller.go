package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ArticleConfig struct {
	*config.Config
}

func New(configuration *config.Config) *ArticleConfig {
	return &ArticleConfig{configuration}
}

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
	imageURL, err := uploadToS3(file, filename)
	utils.HandleError(err, "Error uploading image to S3")

	// Créer un nouvel article dans la base de données
	articleEntry := &dbmodel.ArticleEntry{
		PseudoUser:  user.Pseudo,
		Name:        name,
		Price:       price,
		Size:        size,
		Brand:       brand,
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

func (config *ArticleConfig) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.ArticleDeleteRequest{}
	if errRequest := render.Bind(r, req); errRequest != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}
	errBucket := config.DeleteImageFromS3(req.ArticleId)
	utils.HandleError(errBucket, "Error while deleting image from S3")

	errBDD := config.ArticleRepository.Delete(req.ArticleId)
	utils.HandleError(errBDD, "Error while deleting article from database")

}

func (config *ArticleConfig) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.ArticleUpdateRequest{}
	if errRequest := render.Bind(r, req); errRequest != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}

	article := &dbmodel.ArticleEntry{
		PseudoUser:  req.UserPseudo,
		Name:        req.ArticleName,
		Price:       req.ArticlePrice,
		Size:        req.ArticleSize,
		Brand:       req.ArticleBrand,
		Description: req.ArticleDescription,
		ImageURL:    req.ArticleImage,
	}

	errUpdate := config.ArticleRepository.Update(article)
	utils.HandleError(errUpdate, "Error while updating article in database")

	render.JSON(w, r, map[string]string{"message": "Article updated successfully"})
}

func (config *ArticleConfig) GetArticleID(w http.ResponseWriter, r *http.Request, Id int) (model.ArticleResponse, error) {
	req := &model.ArticleDeleteRequest{}
	article, err := config.ArticleRepository.FindByID(req.ArticleId)
	utils.HandleError(err, "Error while fetching article from database")
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
	return res, nil
}
