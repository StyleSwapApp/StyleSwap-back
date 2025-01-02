package database

import (
	"StyleSwap/database/dbmodel"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("StyleSwap.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected and migrated")
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&dbmodel.ArticleEntry{},
		&dbmodel.UserEntry{},
		&dbmodel.Messages{},
	)
	log.Println("Database migrated successfully")
}
