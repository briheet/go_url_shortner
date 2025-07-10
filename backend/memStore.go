package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type urlStore interface {
	Add(entry *Url) error
	GetByID(urlID string) (Url, error)
	GetByShortURL(shortUrl string) (Url, error)
	Update(urlId string, entry *Url) error
	List(userToken string) ([]Url, error)
	Remove(urlId string) error
}

type userStore interface {
	Add(email, hashedPassword string) (*User, error)
	GetById(userToken string) (User, error)
	GetByEmail(email string) (User, error)
	Update(userToken string, entry *User) error
	Remove(userToken string) error
}

type refreshTokenStore interface {
	GenerateRefreshToken(userId string, ttl time.Duration) (*RefreshToken, error)
	GetRefreshToken(tokenString string) (*RefreshToken, error)
	RevokeRefreshToken(token *RefreshToken) error
}

type urlStoreImpl struct {
	db *gorm.DB
}

type userStoreImpl struct {
	db *gorm.DB
}

type refreshTokenStoreImpl struct {
	db *gorm.DB
}

func (s *urlStoreImpl) Add(entry *Url) error {
	result := s.db.Create(entry)
	return result.Error
}

func (s *urlStoreImpl) GetByID(urlID string) (Url, error) {
	var entry Url
	result := s.db.First(&entry, "id = ?", urlID)
	return entry, result.Error
}

func (s *urlStoreImpl) GetByShortURL(shortUrl string) (Url, error) {
	var entry Url
	result := s.db.First(&entry, "short_url = ?", shortUrl)
	return entry, result.Error
}

func (s *urlStoreImpl) Update(shortUrl string, entry *Url) error {
	result := s.db.Save(entry)
	return result.Error
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

func (s *userStoreImpl) Add(email, hashedPassword string) (*User, error) {
	entry := &User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hashedPassword,
	}
	result := s.db.Create(entry)
	return entry, result.Error
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

func (s *refreshTokenStoreImpl) GenerateRefreshToken(userId string, ttl time.Duration) (*RefreshToken, error) {
	tokenId := uuid.NewString()
	expiresAt := time.Now().Add(ttl)

	refreshToken := &RefreshToken{
		UserId:    userId,
		Token:     tokenId,
		ExpiresAt: expiresAt,
	}

	if err := s.db.Create(refreshToken).Error; err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *refreshTokenStoreImpl) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	result := s.db.First(&refreshToken, "token = ?", tokenString)

	if result.Error != nil {
		return nil, result.Error
	}

	return &refreshToken, nil
}

func (s *refreshTokenStoreImpl) RevokeRefreshToken(token *RefreshToken) error {
	result := s.db.Save(token)
	return result.Error
}

func initDB() (*gorm.DB, error) {
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

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func VerifyPassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}

func GetUserIDFromCtx(r *http.Request) string {
	userID := r.Context().Value("userID").(string)
	return userID
}
