package main

import (
	"log"
	"net/http"
)

func main() {
	serverHandler := http.NewServeMux()
	server := http.Server{
		Handler: serverHandler,
		Addr:    ":8080",
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("the server is not working: %v", err)
	}
}
