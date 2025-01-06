package dbmodel

import (
	"errors"
	"fmt"
	"time"
	"gorm.io/gorm"
)

type UserEntry struct {
	gorm.Model
	FName     string    `json:"user_first_name"`
	LName     string    `json:"user_last_name"`
	Civilite  string    `json:"user_civilite"`
	Address   string    `json:"user_address"`
	City      string    `json:"user_city"`
	Country   string    `json:"user_country"`
	UserEmail string    `json:"user_email" gorm:"unique"`	
	Password  string    `json:"user_password"`
	Pseudo    string    `json:"pseudo" gorm:"unique"`
	BirthDate time.Time `json:"user_birthdate"`
}

type UserRepository interface {
	Create(entry *UserEntry) error
	FindAll() ([]UserEntry, error)
	FindByID(id int) (*UserEntry, error)
	Update(id int, entry *UserEntry) error
	Delete(id int) error
	FindByEmail(email string) (*UserEntry, error)
	FindByPseudo(pseudo string) (UserEntry, error)
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
	if err := r.db.Where("id = ? ", id).First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *userRepository) Update(id int, updatedData *UserEntry) error {
	// Chercher l'utilisateur existant par ID
	var existingUser UserEntry
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return err
	}

	// Comparer chaque champ pour voir s'il a changé et mettre à jour seulement ceux qui ont changé
	userModif := existingUser
	if existingUser.FName != updatedData.FName {
		userModif.FName = updatedData.FName
	}
	if existingUser.LName != updatedData.LName {
		userModif.LName = updatedData.LName
	}
	if existingUser.UserEmail != updatedData.UserEmail {
		allUsers, err := r.FindAll()
		if err != nil {
			return fmt.Errorf("error fetching all users: %w", err)
		}
		for _, entry := range allUsers {
			if entry.UserEmail == updatedData.UserEmail {
				return fmt.Errorf("user with email %s already exists", updatedData.UserEmail)
			}
		}
		userModif.UserEmail = updatedData.UserEmail
	}
	if existingUser.Civilite != updatedData.Civilite {
		userModif.Civilite = updatedData.Civilite
	}
	if existingUser.Address != updatedData.Address {
		userModif.Address = updatedData.Address
	}
	if existingUser.City != updatedData.City {
		userModif.City = updatedData.City
	}
	if existingUser.Country != updatedData.Country {
		userModif.Country = updatedData.Country
	}
	if existingUser.Password != updatedData.Password {
		return fmt.Errorf("password cannot be updated")
	}
	if existingUser.Pseudo != updatedData.Pseudo {
		allUsers, err := r.FindAll()
		if err != nil {
			return fmt.Errorf("error fetching all users: %w", err)
		}
		for _, entry := range allUsers {
			if entry.Pseudo == updatedData.Pseudo {
				return fmt.Errorf("user with pseudo %s already exists", updatedData.Pseudo)
			}
		}
		userModif.Pseudo = updatedData.Pseudo
	}
	if existingUser.BirthDate != updatedData.BirthDate {
		userModif.BirthDate = updatedData.BirthDate
	}

	// Mettre à jour les champs modifiés
	if err := r.db.Model(&existingUser).Updates(userModif).Error; err != nil {
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

func (r *userRepository) FindByPseudo(pseudo string) (UserEntry, error) {
	var entry UserEntry
	if err := r.db.Where("pseudo = ?", pseudo).First(&entry).Error; err != nil {
		return entry, err
	}
	return entry, nil
}

func (r *userRepository) FindByEmail(email string) (*UserEntry, error) {
	var entry UserEntry
	if err := r.db.Where("user_email = ?", email).First(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}
