package main

import (
	"fmt"
	"net/http"
)

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resourcePath := req.PathValue("route")
	method := req.Method
	switch {
	case method == http.MethodGet:
		switch {
		case usersPathRegEx.MatchString(resourcePath):
			h.ListUsers(w, req)
			return
		case urlsPathRegEx.MatchString(resourcePath):
			h.ListUrls(w, req)
			return
		}
	case method == http.MethodGet:
		switch {
		case usersPathRegEx.MatchString(resourcePath):
			h.GetUser(w, req)
			return
		case urlsPathRegEx.MatchString(resourcePath):
			h.GetUrl(w, req)
			return
		}
	case method == http.MethodPost:
		switch {
		case usersPathRegEx.MatchString(resourcePath):
			h.CreateUser(w, req)
			return
		case urlsPathRegEx.MatchString(resourcePath):
			h.CreateUrl(w, req)
			return
		}
	case method == http.MethodDelete:
		switch {
		case usersPathRegEx.MatchString(resourcePath):
			h.DeleteUser(w, req)
			return
		case urlsPathRegEx.MatchString(resourcePath):
			h.DeleteUrl(w, req)
			return
		}
	case method == http.MethodPut:
		switch {
		case usersPathRegEx.MatchString(resourcePath):
			h.UpdateUser(w, req)
			return
		case urlsPathRegEx.MatchString(resourcePath):
			h.UpdateUrl(w, req)
			return
		}
	}
}

func (h *apiHandler) CreateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url created", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) ListUrls(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Urls Listed", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) GetUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url got", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) UpdateUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url updated", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) DeleteUrl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Url deleted", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User created", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) ListUsers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Users listed", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User got", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User updated", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}

func (h *apiHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("User deleted", req.PathValue("route"))
	fmt.Fprintln(w, "API endpoint reached", req.URL.Path)
}
