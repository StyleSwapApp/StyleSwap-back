package register

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type UserConfig struct {
	*config.Config
}

func New(configuration *config.Config) *UserConfig {
	return &UserConfig{configuration}
}

func (config *UserConfig) UserHandler(w http.ResponseWriter, r *http.Request) {
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}

	userAll, err := config.UserRepository.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Vérifier si le pseudo existe déjà
	for _, entry := range userAll {
		if entry.Pseudo == req.Pseudo{
			render.JSON(w, r, map[string]string{"error": "Pseudo already exists"})
			return
		}
	}

	dateB, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid date format"})
		return
	}

	HASH := hashedPassword(req.UserPW)
	userEntry := &dbmodel.UserEntry{
		FName:     req.UserFName,
		LName:     req.UserLName,
		Civilite:  req.Civilite,
		Address:   req.Address,
		City:      req.City,
		Country:   req.Country,
		UserEmail: req.UserEmail,
		Password:  HASH,
		Pseudo:    req.Pseudo,
		BirthDate: dateB,
	}

	config.UserRepository.Create(userEntry)
	token, err := auth.GenerateToken("StyleSwap", req.UserEmail)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func hashedPassword(password string) string {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bcryptPassword)
}