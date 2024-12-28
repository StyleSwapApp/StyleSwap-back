package dbmodel

import (
	"errors"

	"gorm.io/gorm"
)

type ArticleEntry struct {
	gorm.Model
	PseudoUser 	string		`json:"user_pseudo"`
	Name 		string  	`json:"article_name"`
	Price 		int 		`json:"article_price"`
	Description string  	`json:"article_description"`
	ImageURL    string  	`json:"article_image"`
}

type ArticleRepository interface {
	Create(entry *ArticleEntry) error
	FindAll() ([]ArticleEntry, error)
	FindByID(id int) (*ArticleEntry, error)
	FindByPseudo(pseudo string) ([]ArticleEntry, error)
	Delete(id int) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleEntryRepository(db *gorm.DB) *articleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(entry *ArticleEntry) error {
	if entry.PseudoUser == "" {
		return errors.New("missing required IDUser fields")
	}
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

func (r *articleRepository) FindByPseudo(pseudo string) ([]ArticleEntry, error) {
	var entries []ArticleEntry
	if err := r.db.Where("pseudo = ?", pseudo).Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *articleRepository) Delete(id int) error {
	if err := r.db.Delete(&ArticleEntry{}, id).Error; err != nil {
		return err
	}
	return nil
}