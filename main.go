package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	serverHandler := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: serverHandler,
	}

	serverHandler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serverHandler.HandleFunc("/healthz", handlerReadiness)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
