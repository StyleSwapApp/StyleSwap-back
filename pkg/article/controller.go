package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
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
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Unable to parse multipart form"})
		return
	}

	username := r.FormValue("userPseudo")
	user, err := config.UserRepository.FindByPseudo(username)
	if err != nil {	
		render.JSON(w, r, map[string]string{"error": "User not found"})
		return
	}
	if user == "" {
		render.JSON(w, r, map[string]string{"error": "User not found"})
		return
	}
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid price format"})
		return
	}
	description := r.FormValue("description")

	// Vérification que tous les champs sont fournis
	if name == "" || price == 0 || description == "" {
		render.JSON(w, r, map[string]string{"error": "Missing required fields"})
		return
	}

	// Récupérer le fichier de la requête (champ "image")
	file, _, err := r.FormFile("image")
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Unable to get file from form"})
		return
	}
	defer file.Close()

	// Générer un nom unique pour le fichier sur S3
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s.png", uniqueID) // Utiliser un nom unique pour éviter les collisions

	// Télécharger l'image sur S3
	imageURL, err := uploadToS3(file, filename)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to upload image to S3"})
		return
	}

	// Créer un nouvel article dans la base de données
	articleEntry := &dbmodel.ArticleEntry{
		PseudoUser:		user,
		Name:			name,
		Price:      	price,
		Description:	description,
		ImageURL:    	imageURL,
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
	if Pseudo != "" {
		articles, err = config.ArticleRepository.FindByPseudo(Pseudo)
	} else {
		articles, err = config.ArticleRepository.FindAll()
	}

	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Error while fetching articles from the database"})
		return
	}

	for _, article := range articles {
		res := model.ArticleResponse{
			UserPseudo:  			article.PseudoUser,
			ArticleName:			article.Name,
			ArticlePrice:       	article.Price,
			ArticleDescription: 	article.Description,
			ArticleImage:    		article.ImageURL,
		}
		render.JSON(w, r, res)
	}
}

func (config *ArticleConfig) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.ArticleDeleteRequest{}
	if errRequest := render.Bind(r, req); errRequest != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}
	errBucket := config.DeleteImageFromS3(req.ImageURL)

	if errBucket != nil {
		render.JSON(w, r, map[string]string{"error": "Error while deleting image from S3"})
		return
	}

	errBDD := config.ArticleRepository.Delete(req.ArtileId)

	if errBDD != nil {
		render.JSON(w, r, map[string]string{"error": "Error while deleting article from the database"})
		return
	}
}