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

type ArticleRepository interface {
	Create(entry *ArticleEntry) error
	FindAll() ([]ArticleEntry, error)
	FindByID(id int) (*ArticleEntry, error)
}

func NewArticleRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) Create(entry *ArticleEntry) error {
	if err := r.db.Create(entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) FindAll() ([]ArticleEntry, error) {
	var entries []ArticleEntry
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *articleRepository) FindByID(id int) (*ArticleEntry, error) {
	var entry ArticleEntry
	if err := r.db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}