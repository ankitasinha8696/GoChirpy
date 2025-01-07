package main

import (
	"log"
	"net/http"
    "sync/atomic"
    "fmt"
    "encoding/json"
    "strings"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    // Write the Content-Type header
    w.Header().Set("Content-Type", "text/html")

    // Write the status code
    w.WriteHeader(http.StatusOK)

    // Write the body text
    w.Write([]byte("OK"))

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

type Chirp struct {
    Body string `json:"body"`
}

type Response struct {
    Valid bool   `json:"valid"`
    Cleaned_Body  string `json:"cleaned_body"`
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
    fmt.Fprintf(w, `<html>
                        <body>
                            <h1>Welcome, Chirpy Admin</h1>
                            <p>Chirpy has been visited %d times!</p>
                        </body>
                    </html>`, hits)
}

// Implement the resetHandler method
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits.Store(0)
    fmt.Fprintln(w, "Hits counter reset to 0")
}

// Implement the validateChirp method
func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {
    var chirp Chirp
    err := json.NewDecoder(r.Body).Decode(&chirp)
    if err != nil {
        http.Error(w, `{"error": "Something went wrong"}`, 400)
        return
    }

    if len(chirp.Body) > 140 {
        http.Error(w, `{"error": "Chirp is too long"}`, 400)
        return
    }

    profane_words := [9]string{" kerfuffle ", " sharbert ", " fornax ", " Kerfuffle ", " Sharbert ", " Fornax ", " KERFUFFLE ", " SHARBERT ", " FORNAX "}
    stringListWithoutProfanity := []string{}
    chirpString := chirp.Body

    for _, word := range profane_words {
        startIndex := strings.Index(chirpString, word)
        if startIndex != -1 {
                stringListWithoutProfanity = strings.Split(chirpString, word)
                chirpString = strings.Join(stringListWithoutProfanity, " **** ")
                //chirpString = strings.Replace(strings.ToLower(chirpString), word, " **** ", 1)
        }
    }

    chirp.Body = chirpString

    response := Response{
        Valid: true,
        Cleaned_Body:  chirp.Body,
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, `{"error": "Something went wrong"}`, 400)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write(jsonResponse)
}

func main() {

    newMux := http.NewServeMux()

    apiConfig := &apiConfig{}

    newMux.HandleFunc("GET /api/healthz", readinessHandler)

    fileServer := http.FileServer(http.Dir("."))

    newMux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

    // Register the metricsHandler on the /metrics path
    newMux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)

    // Register the resetHandler on the /reset path
    newMux.HandleFunc("POST /admin/reset", apiConfig.resetHandler)

    // Register the resetHandler on the /validate_chirp path
    newMux.HandleFunc("POST /api/validate_chirp", apiConfig.validateChirp)

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



