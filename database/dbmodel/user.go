package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type UserEntry struct {
	gorm.Model
	FName     string `json:"user_first_name"`
	LName     string `json:"user_last_name"`
	UserEmail string `json:"user_email"`
	Password  string `json:"user_password"`
	BirthDate time.Time `json:"user_birthdate"`
}

type UserRepository interface {
	Create(entry *UserEntry) error
	FindAll() ([]UserEntry, error)
	FindByID(id int) (*UserEntry, error)
	Update(entry *UserEntry) error
	Delete(id int) error
	FindByEmail(email string) (*UserEntry, error)
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

func (r *userRepository) Update(entry *UserEntry) error {
	if err := r.db.Save(entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(id int) error {
	if err := r.db.Delete(&UserEntry{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByEmail(email string) (*UserEntry, error) {
	var entry UserEntry
	if err := r.db.Where("user_email = ?", email).First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

