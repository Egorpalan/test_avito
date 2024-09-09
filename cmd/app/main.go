package main

import (
	"log"
	"net/http"
	"os"

	"tender_service/internal/db"
	"tender_service/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
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

	r.Get("/api/tenders", handlers.GetTendersHandler)
	r.Post("/api/tenders/new", handlers.CreateTenderHandler)
    r.Get("/api/tenders/my", handlers.GetUserTendersHandler)
    r.Patch("/api/tenders/{tenderId}/edit", handlers.UpdateTenderHandler)
    r.Put("/api/tenders/{tenderId}/rollback/{version}", handlers.RollbackTenderHandler)

    r.Post("/api/bids/new", handlers.CreateBidHandler)
    r.Get("/api/bids/my", handlers.GetUserBidsHandler)
	r.Get("/api/bids/{tenderId}/list", handlers.GetBidsForTenderHandler)
	r.Patch("/api/bids/{bidId}/edit", handlers.UpdateBidHandler)
    r.Get("/api/bids/{tenderId}/reviews", handlers.GetReviewsHandler) 

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = "0.0.0.0:8080"
	}

	log.Printf("Server is running on %s\n", address)
	log.Fatal(http.ListenAndServe(address, r))
}
