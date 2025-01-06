package dbmodel

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ArticleEntry struct {
	gorm.Model
	PseudoUser  string `json:"user_pseudo"`
	Name        string `json:"article_name"`
	Price       int    `json:"article_price"`
	Size		string `json:"article_size"`
	Brand       string `json:"article_brand"`
	Description string `json:"article_description"`
	ImageURL    string `json:"article_image"`
}

type ArticleRepository interface {
	Create(entry *ArticleEntry) error
	FindAll() ([]ArticleEntry, error)
	FindByID(id int) (*ArticleEntry, error)
	FindByPseudo(pseudo string) ([]ArticleEntry, error)
	Delete(id int) error
	Update(entry *ArticleEntry) error
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
	if err := r.db.Where("pseudo_user = ?", pseudo).Find(&entries).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}
	return entries, nil
}

func (r *articleRepository) Delete(id int) error {
    if err := r.db.Where("id = ?", id).Delete(&ArticleEntry{}).Error; err != nil {
        return err
    }
    return nil
}

func (r *articleRepository) Update(entry *ArticleEntry) error {
	if entry.ID == 0 {
		return errors.New("missing required IDUser fields")
	}
	// Vérifier si l'entrée existe déjà dans la base de données
	var existingEntry ArticleEntry
	if err := r.db.First(&existingEntry, entry.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// L'enregistrement n'existe pas
			return fmt.Errorf("entry with ID %v not found", entry.ID)
		}
    // Autre erreur liée à la base de données
    	return err
	}

	// Si l'enregistrement existe, effectuer la sauvegarde
	if err := r.db.Save(entry).Error; err != nil {
		return err
	}
	return nil
}
