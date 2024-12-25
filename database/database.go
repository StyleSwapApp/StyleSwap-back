package database

import (
	"StyleSwap/database/dbmodel"
	"fmt"
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
    var err error

    // Lire les variables d'environnement
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")

    // Construire l'URL de connexion (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)

    // Ouvrir la connexion à la base de données MariaDB
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Appliquer les migrations
    DB.AutoMigrate(&dbmodel.ArticleEntry{})
    log.Println("Database connected and migrated")
}

func Migrate(db *gorm.DB) {
    db.AutoMigrate(
        &dbmodel.ArticleEntry{},
    )
    log.Println("Database migrated successfully")
}