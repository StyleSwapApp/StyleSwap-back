package login

import (
	"StyleSwap/config"
	"StyleSwap/pkg/auth"
	"StyleSwap/pkg/model"
	"StyleSwap/utils"
	"encoding/json"
	"net/http"

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