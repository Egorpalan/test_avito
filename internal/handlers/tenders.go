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

func GetTendersHandler(w http.ResponseWriter, r *http.Request) {
	serviceType := r.URL.Query().Get("serviceType")

	tenders, err := models.GetAllTenders(db.DB, serviceType)
	if err != nil {
		log.Printf("Error getting tenders: %v", err)
		http.Error(w, "Failed to get tenders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}

func CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var tender models.Tender
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenderID, err := createTenderInDB(&tender)
	if err != nil {
		log.Printf("Failed to create tender: %v", err)
		http.Error(w, "Failed to create tender", http.StatusInternalServerError)
		return
	}

	tender.ID = tenderID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tender); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func createTenderInDB(tender *models.Tender) (int, error) {
	var id int
	err := db.DB.QueryRow(`
		INSERT INTO tenders (name, description, service_type, status, organization_id, creator_username, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id
	`, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.CreatorUsername).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetUserTendersHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username query parameter is required", http.StatusBadRequest)
		return
	}

	tenders, err := models.GetUserTenders(db.DB, username)
	if err != nil {
		log.Printf("Error getting user tenders: %v", err)
		http.Error(w, "Failed to get user tenders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdateTenderHandler(w http.ResponseWriter, r *http.Request) {
	tenderID, err := strconv.Atoi(chi.URLParam(r, "tenderId"))
	if err != nil {
		http.Error(w, "tenderID query parameter is required", http.StatusBadRequest)
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

	tender, err := models.UpdateTender(db.DB, tenderID, requestBody.Name, requestBody.Description)
	if err != nil {
		log.Printf("Error updating tender: %v", err)
		http.Error(w, "Failed to update tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tender); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func RollbackTenderHandler(w http.ResponseWriter, r *http.Request) {
	tenderID, err := strconv.Atoi(chi.URLParam(r, "tenderId"))
	if err != nil {
		http.Error(w, "Invalid tenderID", http.StatusBadRequest)
		return
	}

	version, err := strconv.Atoi(chi.URLParam(r, "version"))
	if err != nil {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}

	tender, err := models.RollbackTender(db.DB, tenderID, version)
	if err != nil {
		log.Printf("Error rolling back tender: %v", err)
		http.Error(w, "Failed to rollback tender", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tender); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
