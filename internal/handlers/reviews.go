package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tender_service/internal/db"
	"tender_service/internal/models"

	"github.com/go-chi/chi/v5"
)

func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	tenderID, err := strconv.Atoi(chi.URLParam(r, "tenderId"))
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	authorUsername := r.URL.Query().Get("authorUsername")
	organizationID, err := strconv.Atoi(r.URL.Query().Get("organizationId"))
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	reviews, err := models.GetReviewsForBid(db.DB, tenderID, authorUsername, organizationID)
	if err != nil {
		log.Printf("Error getting reviews: %v", err)
		http.Error(w, "Failed to get reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
