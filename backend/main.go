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
	usersPathRegEx       = regexp.MustCompile(`^users\/*$`)
	usersPathWithIdRegEx = regexp.MustCompile(`^users\/([a-z0-9-]+)$`)
	urlsPathRegEx        = regexp.MustCompile(`^urls\/*$`)
	urlsPathWithIdRegEx  = regexp.MustCompile(`^urls\/([a-z0-9]+)$`)
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

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(req.Context(), "userID", userID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func main() {
	db, dbErr := initDB()

	if dbErr != nil {
		log.Fatalln("Error initializing database:", dbErr)
	}

	db.AutoMigrate(&User{}, &Url{}, &RefreshToken{})

	userStoreImpl := &userStoreImpl{db: db}

	authService := &authServiceImpl{
		userDb:          userStoreImpl,
		refreshTokenDb:  &refreshTokenStoreImpl{db: db},
		jwtSecret:       []byte(os.Getenv("JWT_SECRET")),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 24 * time.Hour,
	}

	authHandler := &authHandler{
		authService: authService,
	}

	apiHandler := &apiHandler{
		urlDb:  &urlStoreImpl{db: db},
		userDb: userStoreImpl,
	}

	http.Handle("/{short_url}", &shortUrlHandler{
		urlDb: &urlStoreImpl{db: db},
	})
	http.HandleFunc("POST /api/auth/register", authHandler.RegisterUser)
	http.HandleFunc("POST /api/auth/login", authHandler.LoginUser)
	http.HandleFunc("POST /api/auth/refresh", authHandler.RefreshToken)
	http.Handle("/api/{route...}", authMiddleware(authService)(apiHandler))

	log.Println("Starting application on port", 8090)
	err := http.ListenAndServe(":8090", nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("Server Closed")
	} else if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
