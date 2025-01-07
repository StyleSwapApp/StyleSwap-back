package login

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
	"golang.org/x/crypto/bcrypt"
)

type UserConfig struct {
	*config.Config
}

func New(configuration *config.Config) *UserConfig {
	return &UserConfig{configuration}
}

// LoginHandler use to login a user

func (config *UserConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.LoginRequest{}
	if err := req.Bind(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginEntry, err := config.UserRepository.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, entry := range loginEntry {
		if entry.Pseudo == req.Pseudo || entry.UserEmail == req.UserEmail {
			if bcrypt.CompareHashAndPassword([]byte(entry.Password), []byte(req.UserPW)) == nil {
				token, errGenerateToken := auth.GenerateToken("StyleSwap", entry.Pseudo)
				utils.HandleError(errGenerateToken, "Error while generating token")
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"token": token})
				return
			}
		}
	}
	http.Error(w, "Invalid Email/Pseudo or password", http.StatusUnauthorized)
}

// UpdateHandler use to update a user

func (config *UserConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}

	dateB, err := time.Parse("2006-01-02", req.BirthDate)
	utils.HandleError(err, "Error while parsing date")

	userEntry := &dbmodel.UserEntry{
		FName:     req.UserFName,
		LName:     req.UserLName,
		Civilite:  req.Civilite,
		Address:   req.Address,
		City:      req.City,
		Country:   req.Country,
		BirthDate: dateB,
	}

	config.UserRepository.Update(req.UserID, userEntry)
	render.JSON(w, r, "User updated")
}

// DeleteHandler use to delete a user

func (config *UserConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.UserSearchRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}

	config.UserRepository.Delete(req.UserID)
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

func (config *UserConfig) NewPassword(w http.ResponseWriter, r *http.Request, Id int) error {
	return nil
}
