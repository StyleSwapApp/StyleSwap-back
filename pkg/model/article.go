package model

import (
	"errors"
	"net/http"
)

type ArticleDeleteRequest struct {
	ArtileId int `json:"article_id"`
	ImageURL string `json:"image_url"`
}

func (a *ArticleDeleteRequest) Bind(r *http.Request) error {
	if a.ArtileId == 0 {
		return errors.New("missing required ArticleId fields")
	}
	return nil
}

type ArticleResponse struct {
	UserPseudo         string `json:"user_pseudo"`
	ArticleName        string `json:"article_name"`
	ArticlePrice       int    `json:"article_price"`
	ArticleDescription string `json:"article_description"`
	ArticleImage       string `json:"article_image"`
}
