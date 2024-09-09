package models

import (
	"database/sql"
	"time"
)

type Bid struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	TenderID        int       `json:"tenderId"`
	OrganizationID  int       `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
	CreatedAt       time.Time `json:"createdAt"`
}

func CreateBid(db *sql.DB, name, description, status string, tenderID, organizationID int, creatorUsername string) (*Bid, error) {
	query := `
        INSERT INTO bids (name, description, status, tender_id, organization_id, creator_username, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        RETURNING id, created_at
    `
	bid := &Bid{}
	err := db.QueryRow(query, name, description, status, tenderID, organizationID, creatorUsername).Scan(
		&bid.ID,
		&bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	bid.Name = name
	bid.Description = description
	bid.Status = status
	bid.TenderID = tenderID
	bid.OrganizationID = organizationID
	bid.CreatorUsername = creatorUsername

	return bid, nil
}

func GetBidsByUsername(db *sql.DB, username string) ([]Bid, error) {
	var bids []Bid
	query := `SELECT id, name, description, status, tender_id, organization_id, creator_username, created_at
              FROM bids WHERE creator_username = $1`
	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid Bid
		if err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.OrganizationID, &bid.CreatorUsername, &bid.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func GetBidsByTenderID(db *sql.DB, tenderID int) ([]Bid, error) {
	var bids []Bid
	query := `SELECT id, name, description, status, tender_id, organization_id, creator_username, created_at
              FROM bids WHERE tender_id = $1`
	rows, err := db.Query(query, tenderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bid Bid
		if err := rows.Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.OrganizationID, &bid.CreatorUsername, &bid.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func UpdateBid(db *sql.DB, bidID int, name string, description string) (Bid, error) {
	var bid Bid
	query := `UPDATE bids SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description, status, tender_id, organization_id, creator_username, created_at`
	err := db.QueryRow(query, name, description, bidID).Scan(&bid.ID, &bid.Name, &bid.Description, &bid.Status, &bid.TenderID, &bid.OrganizationID, &bid.CreatorUsername, &bid.CreatedAt)
	if err != nil {
		return Bid{}, err
	}

	return bid, nil
}