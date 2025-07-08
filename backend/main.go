package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type apiHandler struct {
	urlDb  urlStore
	userDb userStore
}
type shortUrlHandler struct {
	urlDb urlStore
}

var (
	usersPathRegEx       = regexp.MustCompile(`^users\/*$`)
	usersPathWithIdRegEx = regexp.MustCompile(`^users\/([a-z0-9-]+)$`)
	urlsPathRegEx        = regexp.MustCompile(`^urls\/*$`)
	urlsPathWithIdRegEx  = regexp.MustCompile(`^urls\/([a-z0-9]+)$`)
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		userToken := req.Header.Get("Authorization")

		if userToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func main() {
	db, dbErr := initDB()

	if dbErr != nil {
		fmt.Println("Error initializing database:", dbErr)
	}

	db.AutoMigrate(&User{}, &Url{})

	apiHandler := &apiHandler{
		urlDb:  &urlStoreImpl{db: db},
		userDb: &userStoreImpl{db: db},
	}

	http.Handle("/{short_url}", &shortUrlHandler{
		urlDb: &urlStoreImpl{db: db},
	})
	http.HandleFunc("POST /api/register", apiHandler.RegisterUser)
	http.HandleFunc("POST /api/login", apiHandler.LoginUser)
	http.Handle("/api/{route...}", authMiddleware(apiHandler))

	fmt.Println("Starting application on port", 8090)
	err := http.ListenAndServe(":8090", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
