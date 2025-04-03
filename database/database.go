package database

import (
	"log"

	"github.com/TheDummyUser/goRss/config"
	"github.com/TheDummyUser/goRss/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectionDb() (*gorm.DB, error) {
	dsn := config.Envs.GetDSN()
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect db: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&model.User{}, &model.RssData{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("database successfully migrated")
	return db, nil
}
