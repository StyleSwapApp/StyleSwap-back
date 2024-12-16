package dbmodel

import (
	"gorm.io/gorm"
)

type ArticleEntry struct {
	gorm.Model
	Name 		string  `json:"article_name"`
	Price 		int 	`json:"article_price"`
	Description string `json:"article_description"`
}

type articleRepository struct {
	db *gorm.DB
}