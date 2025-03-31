package dbmodel

import (
	"errors"
	"fmt"

	"StyleSwap/pkg/filter"

	"gorm.io/gorm"
)

type ArticleEntry struct {
	gorm.Model
	PseudoUser  string `json:"user_pseudo"`
	Name        string `json:"article_name"`
	Price       int    `json:"article_price"`
	Size        string `json:"article_size"`
	Brand       string `json:"article_brand"`
	Color       string `json:"article_color"`
	Description string `json:"article_description"`
	ImageData   []byte `gorm:"type:longblob" json:"article_image"`
}

type ArticleRepository interface {
	Create(entry *ArticleEntry) error
	FindAll() ([]ArticleEntry, error)
	FindByID(id int) (*ArticleEntry, error)
	FindByPseudo(pseudo string) ([]ArticleEntry, error)
	Delete(id int) error
	Update(entry *ArticleEntry, id int) error
	FindImageByID(Id int) (string, error)
	FindByCriteria(criteria *filter.FilterCriteria) ([]ArticleEntry, error)
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

func (r *articleRepository) FindImageByID(Id int) (string, error) {
	var entry ArticleEntry
	if err := r.db.First(&entry, Id).Error; err != nil {
		return "", err
	}
	return string(entry.ImageData), nil
}

func (r *articleRepository) Delete(id int) error {
	if err := r.db.Where("id = ?", id).Delete(&ArticleEntry{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) Update(entry *ArticleEntry, id int) error {
	if id == 0 {
		return errors.New("missing required ID fields")
	}
	// Vérifier si l'entrée existe déjà dans la base de données
	var existingEntry ArticleEntry
	if err := r.db.First(&existingEntry, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// L'enregistrement n'existe pas
			return fmt.Errorf("entry with ID %v not found", entry.ID)
		}
		// Autre erreur liée à la base de données
		return err
	}

	// Mettre à jour les champs de l'entrée existante
	if existingEntry.PseudoUser != entry.PseudoUser && entry.PseudoUser != "" {
		existingEntry.PseudoUser = entry.PseudoUser
	}
	if existingEntry.Name != entry.Name && entry.Name != "" {
		existingEntry.Name = entry.Name
	}
	if existingEntry.Price != entry.Price && entry.Price != 0 {
		existingEntry.Price = entry.Price
	}
	if existingEntry.Size != entry.Size && entry.Size != "" {
		existingEntry.Size = entry.Size
	}
	if existingEntry.Brand != entry.Brand && entry.Brand != "" {
		existingEntry.Brand = entry.Brand
	}
	if existingEntry.Description != entry.Description && entry.Description != "" {
		existingEntry.Description = entry.Description
	}

	// Si l'enregistrement existe, effectuer la sauvegarde
	if err := r.db.Save(existingEntry).Error; err != nil {
		return err
	}
	return nil
}

func (r *articleRepository) FindByCriteria(criteria *filter.FilterCriteria) ([]ArticleEntry, error) {
	query := "SELECT * FROM article_entries WHERE deleted_at IS NULL"
	args := []interface{}{}

	if criteria.Pseudo != "" {
		query += " AND pseudo_user = ?"
		args = append(args, criteria.Pseudo)
	}
	if criteria.Couleur != "" {
		query += " AND color = ?"
		args = append(args, criteria.Couleur)
	}
	if criteria.Marque != "" {
		query += " AND brand = ?"
		args = append(args, criteria.Marque)
	}
	if criteria.Taille != "" {
		query += " AND size = ?"
		args = append(args, criteria.Taille)
	}

	var articles []ArticleEntry
	if err := r.db.Raw(query, args...).Scan(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}
