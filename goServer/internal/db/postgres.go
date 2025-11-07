package db

import (
	"log"

	"goServer/internal/config"
	"goServer/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.Config) *gorm.DB {
	User := model.User{}
	dial := postgres.Open(cfg.DatabaseURL)
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to conn ect database: %v", err)
	}

	if err := db.AutoMigrate(User); err != nil {
		log.Printf("AutoMigrate warning/error: %v", err)
	}

	return db
}
