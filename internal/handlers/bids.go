// handlers/bids.go
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

func CreateBidHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		Status          string `json:"status"`
		TenderID        int    `json:"tenderId"`
		OrganizationID  int    `json:"organizationId"`
		CreatorUsername string `json:"creatorUsername"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bid, err := models.CreateBid(
		db.DB,
		requestBody.Name,
		requestBody.Description,
		requestBody.Status,
		requestBody.TenderID,
		requestBody.OrganizationID,
		requestBody.CreatorUsername,
	)
	if err != nil {
		log.Printf("Error creating bid: %v", err)
		http.Error(w, "Failed to create bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bid); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetUserBidsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username query parameter is required", http.StatusBadRequest)
		return
	}

	bids, err := models.GetBidsByUsername(db.DB, username)
	if err != nil {
		log.Printf("Error getting bids: %v", err)
		http.Error(w, "Failed to get bids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bids); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetBidsForTenderHandler(w http.ResponseWriter, r *http.Request) {
	tenderID, err := strconv.Atoi(chi.URLParam(r, "tenderId"))
	if err != nil {
		http.Error(w, "tenderId path parameter is required", http.StatusBadRequest)
		return
	}

	bids, err := models.GetBidsByTenderID(db.DB, tenderID)
	if err != nil {
		log.Printf("Error getting bids for tender: %v", err)
		http.Error(w, "Failed to get bids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bids); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateBidHandler(w http.ResponseWriter, r *http.Request) {
	bidID, err := strconv.Atoi(chi.URLParam(r, "bidId"))
	if err != nil {
		http.Error(w, "bidId path parameter is required", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bid, err := models.UpdateBid(db.DB, bidID, requestBody.Name, requestBody.Description)
	if err != nil {
		log.Printf("Error updating bid: %v", err)
		http.Error(w, "Failed to update bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bid); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

