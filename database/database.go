package database

import (
	"StyleSwap/database/dbmodel"
	"log"
	"gorm.io/gorm"
)

// Migrate effectue la migration des modèles
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&dbmodel.ArticleEntry{},
		&dbmodel.UserEntry{},
		&dbmodel.Messages{},
		&dbmodel.Image{},
	)
	log.Println("Database migrated successfully")
}
