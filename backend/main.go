package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type apiHandler struct{}
type shortUrlHandler struct {
	store urlStore
}

var (
	usersPathRegEx       = regexp.MustCompile(`^users\/*$`)
	usersPathWithIdRegEx = regexp.MustCompile(`^users\/([a-z0-9]+)$`)
	urlsPathRegEx        = regexp.MustCompile(`^urls\/*$`)
	urlsPathWithIdRegEx  = regexp.MustCompile(`^urls\/([a-z0-9]+)$`)
)

func (h *shortUrlHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at short URl", req.PathValue("short_url"))
	if req.PathValue("short_url") == "google" {
		http.Redirect(w, req, "https://www.google.com", http.StatusTemporaryRedirect)
	}
}

func main() {
	_, dbErr := initDB()

	if dbErr != nil {
		fmt.Println("Error initializing database:", dbErr)
		os.Exit(1)
	}

	http.Handle("/{short_url}", &shortUrlHandler{})
	http.Handle("/api/{route...}", &apiHandler{})

	fmt.Println("Starting application on port", 8090)
	err := http.ListenAndServe(":8090", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
