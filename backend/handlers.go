package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (h *shortUrlHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url, err := h.urlDb.Get(req.URL.Path[1:])
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

	userId, _ := req.Cookie("user_token")

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	entry := &Url{
		ID:       uuid.NewString(),
		ShortUrl: requestData.ShortUrl,
		LongUrl:  requestData.LongUrl,
		UserId:   userId.Value,
	}

	if err := h.urlDb.Add(entry); err != nil {
		http.Error(w, "Error creating URL", http.StatusInternalServerError)
		fmt.Println("Error creating URL for user", userId.Value, ":", err)
		return
	}

	fmt.Fprintln(w, "Short URL created", requestData.ShortUrl)
}

func (h *apiHandler) ListUrls(w http.ResponseWriter, req *http.Request) {
	userId, _ := req.Cookie("user_token")

	urls, err := h.urlDb.List(userId.Value)
	if err != nil {
		http.Error(w, "Error fetching URLs", http.StatusInternalServerError)
		fmt.Println("Error fetching URLs for user", userId.Value, ":", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, "Error encoding Data", http.StatusInternalServerError)
		return
	}

	fmt.Println("Urls listed for user", userId.Value)
}

func (h *apiHandler) GetUrl(w http.ResponseWriter, req *http.Request) {
	urlId := strings.TrimPrefix(req.PathValue("route"), "urls/")

	entry, err := h.urlDb.Get(urlId)
	if err != nil {
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}

	fmt.Println("Url fetched", req.PathValue("route"))
}

func (h *apiHandler) UpdateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url updated", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) DeleteUrl(w http.ResponseWriter, req *http.Request) {
	urlId := strings.TrimPrefix(req.PathValue("route"), "urls/")

	if err := h.urlDb.Remove(urlId); err != nil {
		http.Error(w, "Error deleting URL", http.StatusInternalServerError)
		fmt.Println("Error deleting URL:", err)
		return
	}

	fmt.Println("Url deleted", req.PathValue("route"))
	fmt.Fprintln(w, "URL deleted", urlId)
}

func (h *apiHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	userId := strings.TrimPrefix(req.PathValue("route"), "users/")

	entry, err := h.userDb.GetById(userId)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entry); err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		fmt.Println("Error encoding user data:", err)
		return
	}

	fmt.Println("User fetched", userId)
}

func (h *apiHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		fmt.Println("Error decoding data:", err)
		return
	}

	userId := strings.TrimPrefix(req.PathValue("route"), "users/")

	hashedPassword, passErr := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if passErr != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user := &User{
		ID:           userId,
		Email:        requestData.Email,
		PasswordHash: string(hashedPassword),
	}

	h.userDb.Update(userId, user)
}

func (h *apiHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	userId := strings.TrimPrefix(req.PathValue("route"), "users/")

	if err := h.userDb.Remove(userId); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		fmt.Println("Error deleting user:", err)
		return
	}

	fmt.Println("User deleted", userId)
	fmt.Fprintln(w, "User deleted", userId)
}
