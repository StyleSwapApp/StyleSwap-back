package usermanagement

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"encoding/json"
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
	userParam := chi.URLParam(r, "id4Update")
	if userParam == "" {
		json.NewEncoder(w).Encode("User ID is required")
		return
	}
	userInt, err := strconv.Atoi(userParam)
	utils.HandleError(err, "Error while converting user ID to integer")

	//vérifier que l'utilisateur existe
	user, err := config.UserRepository.FindByID(userInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	// Vérifier que l'utilisateur est autorisé à modifier le compte
	VerifUser(user, w, r)

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
	if id == "" {
		json.NewEncoder(w).Encode("User ID is required")
		return
	}
	idInt, err := strconv.Atoi(id)
	utils.HandleError(err, "Error while converting user ID to integer")

	user, err := config.UserRepository.FindByID(idInt)
	utils.HandleError(err, "Error while fetching user from database")

	// Vérifier que l'utilisateur est autorisé à supprimer le compte
	VerifUser(user, w, r)

	// Supprimer les articles de l'utilisateur
	articles, err := config.ArticleRepository.FindByPseudo(user.Pseudo)
	utils.HandleError(err, "Error while fetching articles from database")

	for _, article := range articles {
		// Supprimer l'article de la base de données
		errBDD := config.ArticleRepository.Delete(int(article.ID))
		utils.HandleError(errBDD, "Error while deleting article from database")
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

	// Vérifier que l'utilisateur est autorisé à voir les informations
	VerifUser(user, w, r)

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

func VerifUser(user *dbmodel.UserEntry, w http.ResponseWriter, r *http.Request) {
	userConnect, ok := auth.GetUserIDFromContext(r.Context())
	if ok {
		if user.Pseudo != userConnect {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "You are not authorized to delete this user",
			})
			return
		}
	}
}
