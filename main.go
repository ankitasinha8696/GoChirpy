package main

import (
	"log"
	"net/http"
    "sync/atomic"
    "fmt"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    // Write the Content-Type header
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Write the status code
    w.WriteHeader(http.StatusOK)

    // Write the body text
    w.Write([]byte("OK"))

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Safely increment the fileserverHits counter
        cfg.fileserverHits.Add(1)

        // Call the next handler in the chain
        next.ServeHTTP(w, r)
    })
}

// Implement the metricsHandler method
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    hits := cfg.fileserverHits.Load()
    fmt.Fprintf(w, "Hits: %d\n", hits)
}

// Implement the resetHandler method
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    fmt.Fprintln(w, "Hits counter reset to 0")
}

func main() {

    newMux := http.NewServeMux()

    apiConfig := &apiConfig{}

    newMux.HandleFunc("/healthz", readinessHandler)

    fileServer := http.FileServer(http.Dir("."))

    newMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

    // Register the metricsHandler on the /metrics path
    newMux.HandleFunc("/metrics", apiConfig.metricsHandler)

    // Register the resetHandler on the /reset path
    newMux.HandleFunc("/reset", apiConfig.resetHandler)

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



