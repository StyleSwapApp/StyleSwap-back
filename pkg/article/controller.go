package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
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

// ArticleHandler gère la création d'un nouvel article

func (config *ArticleConfig) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse la requête multipart pour obtenir les données de formulaire
	err := r.ParseMultipartForm(10 << 20) // Limite de taille de 10 Mo pour l'image
	utils.HandleError(err, "Error parsing form data")

	// Récupérer le pseudo de l'utilisateur
	User, ok := auth.GetUserIDFromContext(r.Context())
	UserInt, err := strconv.Atoi(User)
	utils.HandleError(err, "Error while converting user ID to integer")
	var userpseudo *dbmodel.UserEntry

	if !ok {
		render.JSON(w, r, map[string]string{"error": "Unauthorized"})
		return
	} else {
		userpseudo, err = config.UserRepository.FindByID(UserInt)
		utils.HandleError(err, "Error while fetching user from the database")
	}

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
	if size == "" || brand == "" || color == "" {
		render.JSON(w, r, map[string]string{"Help": "Vous devriez remplir tous les champs"})
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
		PseudoUser:  userpseudo.Pseudo,
		Name:        name,
		Price:       price,
		Size:        size,
		Brand:       brand,
		Color:       color,
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


// Vérifie que l'utilisateur est autorisé à avoir une action sur l'article 

func VerifArticle(config *ArticleConfig, article *dbmodel.ArticleEntry, w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserIDFromContext(r.Context())
	if ok {
		articleUserId, err := config.UserRepository.FindByPseudo(article.PseudoUser)
		utils.HandleError(err, "Error while fetching article from database")
		userID, err := strconv.Atoi(user)
		utils.HandleError(err, "Error while converting user ID to integer")
		if uint(userID) != articleUserId.ID {
			render.JSON(w, r, map[string]string{"error": "Unauthorized"})
			return
		}
	}
}