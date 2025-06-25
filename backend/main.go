package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type ApiHandler struct{}
type ShortUrlHandler struct{}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	switch {
	case method == http.MethodGet:
		h.ListUrls(w, req)
		return
	case method == http.MethodGet:
		h.GetUrl(w, req)
		return
	case method == http.MethodPost:
		h.CreateUrl(w, req)
		return
	case method == http.MethodDelete:
		h.DeleteUrl(w, req)
		return
	case method == http.MethodPut:
		h.UpdateUrl(w, req)
		return
	}
}

func (h *ApiHandler) CreateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) ListUrls(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) GetUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) UpdateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) DeleteUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /api/", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

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
