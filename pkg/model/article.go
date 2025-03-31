package model

import (
	"errors"
	"net/http"
)

type ArticleDeleteRequest struct {
	ArticleId int    `json:"article_id"`
	ImageURL  string `json:"image_url"`
}

func (a *ArticleDeleteRequest) Bind(r *http.Request) error {
	if a.ArticleId == 0 {
		return errors.New("missing required ArticleId fields")
	}
	return nil
}

type ArticleResponse struct {
	ArticleId          int    `json:"article_id"`
	UserPseudo         string `json:"user_pseudo"`
	ArticleName        string `json:"article_name"`
	ArticlePrice       int    `json:"article_price"`
	ArticleSize        string `json:"article_size"`
	ArticleBrand       string `json:"article_brand"`
	ArticleColor       string `json:"article_color"`
	ArticleDescription string `json:"article_description"`
	ArticleImage       string `json:"article_image"`
}

type ArticleUpdateRequest struct {
	UserPseudo         string `json:"user_pseudo"`
	ArticleName        string `json:"article_name"`
	ArticlePrice       int    `json:"article_price"`
	ArticleSize        string `json:"article_size"`
	ArticleBrand       string `json:"article_brand"`
	ArticleDescription string `json:"article_description"`
	ArticleImage       string `json:"article_image"`
}

func (a *ArticleUpdateRequest) Bind(r *http.Request) error {
	if a.UserPseudo == "" {
		return errors.New("missing required UserPseudo fields")
	}
	if a.ArticleName == "" {
		return errors.New("missing required ArticleName fields")
	}
	if a.ArticlePrice == 0 {
		return errors.New("missing required ArticlePrice fields")
	}
	if a.ArticleDescription == "" {
		return errors.New("missing required ArticleDescription fields")
	}
	if a.ArticleImage == "" {
		return errors.New("missing required ArticleImage fields")
	}
	return nil
}
