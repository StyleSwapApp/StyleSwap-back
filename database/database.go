package database

import (
	"StyleSwap/database/dbmodel"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Variable globale pour la base de données
var DB *gorm.DB

// InitDatabase initialise la connexion à la base de données MySQL
func InitDatabase() {
	// Charger le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Construire le DSN avec les variables d'environnement
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Connecter à la base de données
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}

	log.Println("MySQL Database connected successfully")
}

// Migrate effectue la migration des modèles
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&dbmodel.ArticleEntry{},
		&dbmodel.UserEntry{},
		&dbmodel.Messages{},
	)
	log.Println("Database migrated successfully")
}
