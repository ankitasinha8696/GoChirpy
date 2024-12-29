package main

import (
	"log"
	"net/http"
)

func main() {

    newMux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("."))

    newMux.Handle("/", fileServer)

    server := &http.Server{
        Addr:    ":8080",  // Set the address to listen on port 8080
        Handler: newMux,      // Use the ServeMux as the server's handler
    }

    // Start the server
    log.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }

}



