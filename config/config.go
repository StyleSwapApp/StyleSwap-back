package config

import (
	"StyleSwap/database/dbmodel"
	"StyleSwap/database"
	
    "github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

type Config struct {
	articleRepository   dbmodel.ArticleRepository
}

func New() (*Config, error) {
	config := Config{}
	
	databaseSession, err := gorm.Open(sqlite.Open("StyleSwap_api.db"), &gorm.Config{})
    if err != nil {
        return &config, err
    }

    // Migration des modèles
    database.Migrate(databaseSession)

	// Initialisation des repositories
	config.articleRepository = dbmodel.NewArticleRepository(databaseSession)

	return &config, nil
}