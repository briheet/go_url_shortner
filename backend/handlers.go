package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type urlRequest struct {
	ID       string `json:"id"`
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"long_url"`
}

func (h *shortUrlHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url, err := h.urlDb.GetByShortURL(req.URL.Path[1:])
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, req, url.LongUrl, http.StatusTemporaryRedirect)
}

func (h *authHandler) RegisterUser(w http.ResponseWriter, req *http.Request) {
	var requestData Credentials

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if requestData.Email == "" || requestData.Password == "" {
		http.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(requestData.Email, requestData.Password)
	if err != nil {
		if errors.Is(err, ErrEmailInUse) {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&response)
}

func (h *authHandler) LoginUser(w http.ResponseWriter, req *http.Request) {
	var requestData Credentials

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(requestData.Email, requestData.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := &LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, req *http.Request) {
	var requestData RefreshRequest

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	token, err := h.authService.RefreshAccessToken(requestData.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrExpiredToken) {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := &RefreshResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resourcePath := req.PathValue("route")
	method := req.Method

	switch {
	case method == http.MethodGet:
		switch {
		case urlsPathRegEx.MatchString(resourcePath):
			h.ListUrls(w, req)
			return
		case usersPathWithIdRegEx.MatchString(resourcePath):
			h.GetUser(w, req)
			return
		case urlsPathWithIdRegEx.MatchString(resourcePath):
			h.GetUrl(w, req)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	case method == http.MethodPost:
		switch {
		case urlsPathRegEx.MatchString(resourcePath):
			h.CreateUrl(w, req)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	case method == http.MethodDelete:
		switch {
		case usersPathWithIdRegEx.MatchString(resourcePath):
			h.DeleteUser(w, req)
			return
		case urlsPathWithIdRegEx.MatchString(resourcePath):
			h.DeleteUrl(w, req)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	case method == http.MethodPut:
		switch {
		case usersPathWithIdRegEx.MatchString(resourcePath):
			h.UpdateUser(w, req)
			return
		case urlsPathWithIdRegEx.MatchString(resourcePath):
			h.UpdateUrl(w, req)
			return
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *apiHandler) CreateUrl(w http.ResponseWriter, req *http.Request) {
	var requestData struct {
		ShortUrl string `json:"short_url"`
		LongUrl  string `json:"long_url"`
	}

	userIDFromCtx := GetUserIDFromCtx(req)

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	entry := &Url{
		ID:       uuid.NewString(),
		ShortUrl: requestData.ShortUrl,
		LongUrl:  requestData.LongUrl,
		UserId:   userIDFromCtx,
	}

	if err := h.urlDb.Add(entry); err != nil {
		http.Error(w, "Error creating URL", http.StatusInternalServerError)
		log.Println("Error creating URL for user", userIDFromCtx, ":", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

func (h *apiHandler) ListUrls(w http.ResponseWriter, req *http.Request) {
	userIDFromCtx := GetUserIDFromCtx(req)

	urls, err := h.urlDb.List(userIDFromCtx)
	if err != nil {
		http.Error(w, "Error fetching URLs", http.StatusInternalServerError)
		log.Println("Error fetching URLs for user", userIDFromCtx, ":", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

func (h *apiHandler) GetUrl(w http.ResponseWriter, req *http.Request) {
	urlID := strings.TrimPrefix(req.PathValue("route"), "url/")

	entry, err := h.urlDb.GetByID(urlID)
	if err != nil {
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		log.Println("Error fetching URL :", err)
		return
	}

	userIDFromCtx := GetUserIDFromCtx(req)
	if entry.UserId != userIDFromCtx {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&entry)
}

func (h *apiHandler) UpdateUrl(w http.ResponseWriter, req *http.Request) {
	var requestData struct {
		ShortUrl string `json:"short_url"`
		LongUrl  string `json:"long_url"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		log.Println("Error decoding data :", err)
		return
	}

	urlID := strings.TrimPrefix(req.PathValue("route"), "url/")

	url, err := h.urlDb.GetByID(urlID)
	if err != nil {
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		log.Println("Error fetching URL", urlID, ":", err)
		return
	}

	userID := url.UserId

	userIDFromCtx := GetUserIDFromCtx(req)
	if userIDFromCtx != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	entry := &Url{
		ID:       urlID,
		ShortUrl: requestData.ShortUrl,
		LongUrl:  requestData.LongUrl,
	}

	if err := h.urlDb.Update(urlID, entry); err != nil {
		http.Error(w, "Error updating URL", http.StatusInternalServerError)
		log.Println("Error updating URL", urlID, ":", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (h *apiHandler) DeleteUrl(w http.ResponseWriter, req *http.Request) {
	urlID := strings.TrimPrefix(req.PathValue("route"), "url/")

	url, err := h.urlDb.GetByID(urlID)
	if err != nil {
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		log.Println("Error fetching URL", urlID, ":", err)
		return
	}

	userID := url.UserId

	userIDFromCtx := GetUserIDFromCtx(req)
	if userIDFromCtx != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.urlDb.Remove(urlID); err != nil {
		http.Error(w, "Error deleting URL", http.StatusInternalServerError)
		log.Println("Error deleting URL :", err)
		return
	}

	log.Println("Url deleted", urlID)
	fmt.Fprintln(w, "URL", urlID, "deleted successfully")
}

func (h *apiHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	userID := strings.TrimPrefix(req.PathValue("route"), "user/")

	userIDFromCtx := GetUserIDFromCtx(req)
	if userIDFromCtx != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	entry, err := h.userDb.GetById(userID)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&entry)
}

func (h *apiHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	var requestData Credentials

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		log.Println("Error decoding data :", err)
		return
	}

	userID := strings.TrimPrefix(req.PathValue("route"), "user/")

	userIDFromCtx := GetUserIDFromCtx(req)
	if userIDFromCtx != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hashedPassword, passErr := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if passErr != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	user := &User{
		ID:           userID,
		Email:        requestData.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := h.userDb.Update(userID, user); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		log.Println("Error updating user", userID, ":", err)
		return
	}

	log.Println("User updated:", userID)
	fmt.Fprintln(w, "User", userID, "updated successfully")
}

func (h *apiHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	userID := strings.TrimPrefix(req.PathValue("route"), "user/")

	userIDFromCtx := GetUserIDFromCtx(req)
	if userIDFromCtx != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.userDb.Remove(userID); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		log.Println("Error deleting user", userID, ":", err)
		return
	}

	log.Println("User deleted:", userID)
	fmt.Fprintln(w, "User", userID, "deleted successfully")
}
