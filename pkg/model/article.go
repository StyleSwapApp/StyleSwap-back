package model

import (
	"errors"
	"net/http"
)

type ArticleRequest struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

func (a *ArticleRequest) Bind(r *http.Request) error {
	if a.Name == "" {
		return errors.New("missing required Name field")
	}
	if a.Price == 0 {
		return errors.New("missing required Price field")
	}
	return nil
}