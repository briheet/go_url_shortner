package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type urlStore interface {
	Add(entry *Url) error
	Get(urlId string) (Url, error)
	Update(urlId string, entry *Url) error
	List(userToken string) ([]Url, error)
	Remove(urlId string) error
}

type userStore interface {
	Add(entry *User) error
	GetById(userToken string) (User, error)
	GetByEmail(email string) (User, error)
	Update(userToken string, entry *User) error
	Remove(userToken string) error
}

type urlStoreImpl struct {
	db *gorm.DB
}

type userStoreImpl struct {
	db *gorm.DB
}

func (s *urlStoreImpl) Add(entry *Url) error {
	result := s.db.Create(entry)
	return result.Error
}

func (s *urlStoreImpl) Get(shortUrl string) (Url, error) {
	var entry Url
	result := s.db.First(&entry, shortUrl)
	return entry, result.Error
}

func (s *urlStoreImpl) Update(shortUrl string, entry *Url) error {
	return nil
}

func (s *urlStoreImpl) List(user_id string) ([]Url, error) {
	var urls []Url
	result := s.db.Find(&urls, "user_id = ?", user_id)
	return urls, result.Error
}

func (s *urlStoreImpl) Remove(shortUrl string) error {
	result := s.db.Delete(&Url{
		ShortUrl: shortUrl,
	})
	return result.Error
}

func (s *userStoreImpl) Add(entry *User) error {
	result := s.db.Create(entry)
	return result.Error
}

func (s *userStoreImpl) GetById(userId string) (User, error) {
	var entry User
	result := s.db.First(&entry, "id = ?", userId)
	return entry, result.Error
}

func (s *userStoreImpl) GetByEmail(email string) (User, error) {
	var entry User
	result := s.db.First(&entry, "email = ?", email)
	return entry, result.Error
}

func (s *userStoreImpl) Update(userToken string, entry *User) error {
	result := s.db.Save(entry)
	return result.Error
}

func (s *userStoreImpl) Remove(userId string) error {
	result := s.db.Delete(&User{
		ID: userId,
	})
	return result.Error
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
