package main

import (
	"time"
)

type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Url struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ShortUrl  string    `json:"short_url" gorm:"unique"`
	LongUrl   string    `json:"long_url"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserId"`
}
