package main

import (
	"log"
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    // Write the Content-Type header
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Write the status code
    w.WriteHeader(http.StatusOK)

    // Write the body text
    w.Write([]byte("OK"))

}

func main() {

    newMux := http.NewServeMux()

    newMux.HandleFunc("/healthz", readinessHandler)

    fileServer := http.FileServer(http.Dir("."))

    newMux.Handle("/app/", http.StripPrefix("/app", fileServer))

    server := &http.Server {
        Addr:    ":8080",  // Set the address to listen on port 8080
        Handler: newMux,      // Use the ServeMux as the server's handler
    }

    // Start the server
    log.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }

}



