package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type urlStore interface {
	Add(shortUrl string, entry url) error
	Get(shortUrl string) (url, error)
	Update(shortUrl string, entry url) error
	List() (map[string]url, error)
	Remove(shortUrl string) error
}

type urlStoreImpl struct {
	db *gorm.DB
}

func (s *urlStoreImpl) Add(shortUrl string, entry url) error {
	result := s.db.Create(&entry)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *urlStoreImpl) Get(shortUrl string) (url, error) {
	var entry url
	result := s.db.First(&entry, "short_url = ?", shortUrl)
	if result.Error != nil {
		return url{}, result.Error
	}
	return entry, nil
}

func initDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
