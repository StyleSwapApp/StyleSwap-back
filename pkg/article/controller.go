package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"StyleSwap/utils"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
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

	var userpseudo dbmodel.UserEntry

	if !ok {
		render.JSON(w, r, map[string]string{"error": "Unauthorized"})
		return
	} else {
		userpseudo, err = config.UserRepository.FindByPseudo(User)
		utils.HandleError(err, "Error while fetching user from the database")
	}

	// Extraire les données de l'article du formulaire
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil || price <= 0 {
		render.JSON(w, r, map[string]string{"error": "Invalid price value"})
		return
	}
	size := r.FormValue("size")
	brand := r.FormValue("brand")
	color := r.FormValue("color")
	description := r.FormValue("description")

	// Vérification que tous les champs sont fournis
	if name == "" || price <= 0 || description == "" || size == "" || brand == "" || color == "" {
		render.JSON(w, r, map[string]string{"error": "Missing required fields"})
		return
	}

	// Récupérer le fichier de la requête (champ "image")
	file, _, err := r.FormFile("image")
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Missing image file"})
		return
	}
	defer file.Close()

	// Lire l'image en binaire
	imageData, err := io.ReadAll(file)
	if err != nil {
		utils.HandleError(err, "Erreur lors de la lecture de l'image")
		return
	}

	// Créer un nouvel article dans la base de données
	articleEntry := &dbmodel.ArticleEntry{
		PseudoUser:  userpseudo.Pseudo,
		Name:        name,
		Price:       price,
		Size:        size,
		Brand:       brand,
		Color:       color,
		Description: description,
		ImageData:   imageData,
	}

	// Ajouter l'article à la base de données
	if err := config.ArticleRepository.Create(articleEntry); err != nil {
		render.JSON(w, r, map[string]string{"error": "Error while adding article to the database"})
		return
	}

	// Répondre avec un message de succès
	render.JSON(w, r, map[string]string{"message": "Article added successfully"})
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
