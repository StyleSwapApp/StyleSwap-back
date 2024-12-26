package usermanagement

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
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
		render.JSON(w,r,map[string]string{"error": "Invalid request payload"})
		return 
	}

	dateB, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		render.JSON(w,r,map[string]string{"error": "Invalid date format"})
		return 
	}

	HASH := hashedPassword(req.UserPW)
	userEntry := &dbmodel.UserEntry{
		FName: req.UserFName,
		LName: req.UserLName,
		UserEmail: req.UserEmail,
		Password: HASH,
		Pseudo: req.Pseudo,
		BirthDate: dateB,
	}
	config.UserRepository.Create(userEntry)
}

func hashedPassword(password string) string {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bcryptPassword)
}
