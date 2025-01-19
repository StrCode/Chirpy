package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	serverHandler := http.NewServeMux()
	srv := http.Server{
		Addr:    ":" + port,
		Handler: serverHandler,
	}

	serverHandler.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
