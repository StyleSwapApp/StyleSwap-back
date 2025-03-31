package dbmodel

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	ArticleID uint         `gorm:"index"`                                             // Clé étrangère vers Article
	Article   ArticleEntry `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;"` // Reference to Article struct
	Data      []byte       `json:"Data" gorm:"type:longblob"`
}

// ImageRepository définit les méthodes pour interagir avec le modèle Image
type ImageRepository interface {
	// Créer une nouvelle image
	Create(image *Image) error
}

type imageRepository struct {
	db *gorm.DB
}

// NewImageEntryRepository crée un nouveau repository pour le modèle Image
func NewImageEntryRepository(db *gorm.DB) *imageRepository {
	return &imageRepository{db}
}

func (r *imageRepository) Create(image *Image) error {
	return r.db.Create(image).Error
}
