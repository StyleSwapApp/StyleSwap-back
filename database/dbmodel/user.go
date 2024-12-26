package dbmodel

import (
	"gorm.io/gorm"
)

type UserEntry struct {
	gorm.Model
	FName     string `json:"article_name"`
	LName     string `json:"article_price"`
	UserEmail string `json:"article_description"`
	Password  string `json:"article_image"`
}

type UserRepository interface {
	Create(entry *ArticleEntry) error
	FindAll() ([]ArticleEntry, error)
	FindByID(id int) (*ArticleEntry, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUSerEntryRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(entry *UserEntry) error {
	if err := r.db.Create(entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindAll() ([]UserEntry, error) {
	var entries []UserEntry
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *userRepository) FindByID(id int) (*UserEntry, error) {
	var entry UserEntry
	if err := r.db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}
