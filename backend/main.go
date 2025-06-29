package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var (
	UsersPathRegEx       = regexp.MustCompile(`^users\/*$`)
	UsersPathWithIdRegEx = regexp.MustCompile(`^users\/([a-z0-9]+)$`)
	UrlsPathRegEx        = regexp.MustCompile(`^urls\/*$`)
	UrlsPathWithIdRegEx  = regexp.MustCompile(`^urls\/([a-z0-9]+)$`)
)

func (h *ShortUrlHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at short URl", req.PathValue("short_url"))
	if req.PathValue("short_url") == "google" {
		http.Redirect(w, req, "https://www.google.com", http.StatusTemporaryRedirect)
	}
}

func main() {
	http.Handle("/{short_url}", &ShortUrlHandler{})
	http.Handle("/api/{route...}", &ApiHandler{})

	fmt.Println("Starting application on port", 8090)
	err := http.ListenAndServe(":8090", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
