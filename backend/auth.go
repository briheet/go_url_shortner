package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrEmailInUse         = errors.New("email already in use")
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

type authService interface {
	Register(email, password string) (*User, error)
	GenerateAccessToken(user *User) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	Login(email, password string) (accessToken, refreshToken string, err error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type authServiceImpl struct {
	userDb          userStore
	refreshTokenDb  refreshTokenStore
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func (a *authServiceImpl) Register(email, password string) (*User, error) {
	_, err := a.userDb.GetByEmail(email)
	if err == nil {
		return nil, ErrEmailInUse
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := a.userDb.Add(email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *authServiceImpl) Login(email, password string) (accessToken, refreshTokenString string, err error) {
	user, err := a.userDb.GetByEmail(email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err = a.GenerateAccessToken(&user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.refreshTokenDb.GenerateRefreshToken(user.ID, a.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken.Token, nil
}

func (a *authServiceImpl) GenerateAccessToken(user *User) (string, error) {

	expirationTime := time.Now().Add(a.accessTokenTTL)

	claims := &jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   expirationTime.Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, signErr := token.SignedString(a.jwtSecret)
	if signErr != nil {
		return "", signErr
	}

	return tokenString, nil
}

func (a *authServiceImpl) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (a *authServiceImpl) RefreshAccessToken(refreshTokenString string) (string, error) {
	token, err := a.refreshTokenDb.GetRefreshToken(refreshTokenString)
	if err != nil {
		return "", ErrInvalidToken
	}

	if token.Revoked {
		return "", ErrInvalidToken
	}

	if time.Now().After(token.ExpiresAt) {
		return "", ErrExpiredToken
	}

	user, err := a.userDb.GetById(token.User.ID)
	if err != nil {
		return "", err
	}

	accessToken, err := a.GenerateAccessToken(&user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
