package usermanagement

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserConfig struct {
	*config.Config
}

func New(configuration *config.Config) *UserConfig {
	return &UserConfig{configuration}
}

// UpdateHandler use to update a user

func (config *UserConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "id4Update")
	userInt, err := strconv.Atoi(user)
	utils.HandleError(err, "Error while converting user ID to integer")

	//vérifier que l'utilisateur existe
	_, err = config.UserRepository.FindByID(userInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	req := &model.UserRequest{}
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON format",
		})
		return
	}
	var dateB time.Time
	if req.BirthDate != "" {
		dateB, err = time.Parse("2006-01-02", req.BirthDate)
		utils.HandleError(err, "Error while parsing date")
	}

	userEntry := &dbmodel.UserEntry{
		FName:     req.UserFName,
		LName:     req.UserLName,
		Civilite:  req.Civilite,
		Address:   req.Address,
		City:      req.City,
		Country:   req.Country,
		BirthDate: dateB,
	}

	config.UserRepository.Update(userInt, userEntry)
	render.JSON(w, r, "User updated")
}

// DeleteHandler use to delete a user and their articles

func (config *UserConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id4Delete")
	idInt, err := strconv.Atoi(id)
	utils.HandleError(err, "Error while converting user ID to integer")

	user, err := config.UserRepository.FindByID(idInt)
	utils.HandleError(err, "Error while fetching user from database")

	// Supprimer les articles de l'utilisateur
	articles, err := config.ArticleRepository.FindByPseudo(user.Pseudo)
	utils.HandleError(err, "Error while fetching articles from database")

	for _, article := range articles {
		// Supprimer l'article de la base de données
		errBDD := config.ArticleRepository.Delete(int(article.ID))
		utils.HandleError(errBDD, "Error while deleting article from database")

		// Supprimer l'image de S3
		errBucket := utils.DeleteImageFromS3(article.ImageURL)
		utils.HandleError(errBucket, "Error while deleting image from S3")
	}

	// Supprimer l'utilisateur de la base de données

	errBDD := config.UserRepository.Delete(idInt)
	utils.HandleError(errBDD, "Error while deleting user from database")
	render.JSON(w, r, "User deleted")
}

// GetUserHandler use to get a user

func (config *UserConfig) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	// Récupérer l'ID depuis l'URL
	idParam := chi.URLParam(r, "id")

	// Convertir l'ID en entier
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid user ID",
		})
		return
	}

	// Rechercher l'utilisateur
	user, errUser := config.UserRepository.FindByID(userID)
	if errUser != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	// Construire la réponse
	res := model.UserCompletResponse{
		UserFName: user.FName,
		UserLName: user.LName,
		Civilite:  user.Civilite,
		Address:   user.Address,
		City:      user.City,
		Country:   user.Country,
		UserEmail: user.UserEmail,
		Pseudo:    user.Pseudo,
		BirthDate: user.BirthDate.Format("2006-01-02"),
	}

	// Retourner la réponse JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
