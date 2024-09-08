package main

import (
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/joho/godotenv"
	"tender_service/internal/db"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    err = db.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Get("/api/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    
    address := os.Getenv("SERVER_ADDRESS")
    if address == "" {
        address = "0.0.0.0:8080" 
    }

    log.Printf("Server is running on %s\n", address)
    log.Fatal(http.ListenAndServe(address, r))
}
