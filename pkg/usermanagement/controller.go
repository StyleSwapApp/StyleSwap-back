package usermanagement

import (
	"StyleSwap/config"
	"net/http"
)

type UserConfig struct {
	configuration *config.Config
}

func New(configuration *config.Config) *UserConfig {
	return &UserConfig{configuration}
}

func (config *UserConfig) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}