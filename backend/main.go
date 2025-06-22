package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func headers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /headers")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintln(w, name, h)
		}
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Request received at /hello")
	fmt.Fprintln(w, "hello")
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	err := http.ListenAndServe(":8090", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
