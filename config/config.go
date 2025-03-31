package config

import (
	"StyleSwap/database"
	"StyleSwap/database/dbmodel"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config contient les repositories nécessaires pour l'application
type Config struct {
	// Connexion aux repositories
	ArticleRepository dbmodel.ArticleRepository
	UserRepository    dbmodel.UserRepository
	MessageRepository dbmodel.MessageRepository
	ImageRepository   dbmodel.ImageRepository
}

// New initialise la configuration avec la base de données MySQL
func New() (*Config, error) {
	// Charger les variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur lors du chargement du fichier .env:", err)
	}

	// Construire le DSN MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Se connecter à la base de données MySQL
	databaseSession, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la connexion à la base de données: %v", err)
	}

	// Migration des modèles
	database.Migrate(databaseSession)

	// Initialisation des repositories
	config := &Config{
		ArticleRepository: dbmodel.NewArticleEntryRepository(databaseSession),
		UserRepository:    dbmodel.NewUSerEntryRepository(databaseSession),
		MessageRepository: dbmodel.NewMessageEntryRepository(databaseSession),
		ImageRepository:   dbmodel.NewImageEntryRepository(databaseSession),
	}

	return config, nil
}
