package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type apiHandler struct {
	urlDb  urlStore
	userDb userStore
}
type shortUrlHandler struct {
	urlDb urlStore
}

type authHandler struct {
	authService authService
}

var (
	usersPathRegEx       = regexp.MustCompile(`^user\/*$`)
	usersPathWithIdRegEx = regexp.MustCompile(`^user\/([a-z0-9-]+)$`)
	urlsPathRegEx        = regexp.MustCompile(`^url\/*$`)
	urlsPathWithIdRegEx  = regexp.MustCompile(`^url\/([a-z0-9-]+)$`)
	emailRegEx           = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegEx        = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d@$!%*?&]{8,}$`)
)

func authMiddleware(authService authService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			userIDStr, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			_, err = uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), "userID", userIDStr)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatalln("Error initializing database:", err)
	}

	db.AutoMigrate(&User{}, &Url{}, &RefreshToken{})

	port := os.Getenv("PORT")
	jwtSecretString := os.Getenv("JWT_SECRET")

	userStoreImpl := &userStoreImpl{db: db}
	urlStoreImpl := &urlStoreImpl{db: db}
	authService := &authServiceImpl{
		userDb:          userStoreImpl,
		refreshTokenDb:  &refreshTokenStoreImpl{db: db},
		jwtSecret:       []byte(jwtSecretString),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 24 * time.Hour,
	}

	authHandler := &authHandler{
		authService: authService,
	}
	shortUrlHandler := &shortUrlHandler{
		urlDb: urlStoreImpl,
	}
	apiHandler := &apiHandler{
		urlDb:  urlStoreImpl,
		userDb: userStoreImpl,
	}

	http.Handle("/{short_url}", shortUrlHandler)
	http.HandleFunc("POST /api/auth/register", authHandler.RegisterUser)
	http.HandleFunc("POST /api/auth/login", authHandler.LoginUser)
	http.HandleFunc("POST /api/auth/refresh", authHandler.RefreshToken)
	http.Handle("/api/{route...}", authMiddleware(authService)(apiHandler))

	log.Println("Starting application on port", port)
	err = http.ListenAndServe(":"+port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("Server Closed")
	} else if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
