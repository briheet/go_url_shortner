package main

import (
	"time"
)

type user struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type url struct {
	ShortUrl  string    `json:"short_url" gorm:"primaryKey"`
	LongUrl   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
	User      user      `gorm:"foreignKey:UserID"`
}
