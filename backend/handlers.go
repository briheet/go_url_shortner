package main

import (
	"fmt"
	"net/http"
)

type ApiHandler struct{}
type ShortUrlHandler struct{}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resourcePath := req.PathValue("route")
	method := req.Method
	switch {
	case method == http.MethodGet:
		switch {
		case UsersPathRegEx.MatchString(resourcePath):
			h.ListUsers(w, req)
			return
		case UrlsPathRegEx.MatchString(resourcePath):
			h.ListUrls(w, req)
			return
		}
	case method == http.MethodGet:
		switch {
		case UsersPathRegEx.MatchString(resourcePath):
			h.GetUser(w, req)
			return
		case UrlsPathRegEx.MatchString(resourcePath):
			h.GetUrl(w, req)
			return
		}
	case method == http.MethodPost:
		switch {
		case UsersPathRegEx.MatchString(resourcePath):
			h.CreateUser(w, req)
			return
		case UrlsPathRegEx.MatchString(resourcePath):
			h.CreateUrl(w, req)
			return
		}
	case method == http.MethodDelete:
		switch {
		case UsersPathRegEx.MatchString(resourcePath):
			h.DeleteUser(w, req)
			return
		case UrlsPathRegEx.MatchString(resourcePath):
			h.DeleteUrl(w, req)
			return
		}
	case method == http.MethodPut:
		switch {
		case UsersPathRegEx.MatchString(resourcePath):
			h.UpdateUser(w, req)
			return
		case UrlsPathRegEx.MatchString(resourcePath):
			h.UpdateUrl(w, req)
			return
		}
	}
}

func (h *ApiHandler) CreateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url created", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) ListUrls(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Urls Listed", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) GetUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url got", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) UpdateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url updated", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) DeleteUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url deleted", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User created", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Users listed", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User got", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User updated", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *ApiHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User deleted", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}
