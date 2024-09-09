package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Review struct {
	ID             int       `json:"id"`
	BidID          int       `json:"bidId"`
	AuthorUsername string    `json:"authorUsername"`
	Rating         int       `json:"rating"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"createdAt"`
}

func GetReviewsForBid(db *sql.DB, tenderID int, authorUsername string, organizationID int) ([]Review, error) {
	var reviews []Review
	query := `
        SELECT id, bid_id, author_username, rating, comment, created_at
        FROM reviews
        WHERE bid_id IN (
            SELECT id FROM bids WHERE tender_id = $1 AND organization_id = $2
        ) AND author_username = $3
    `
	rows, err := db.Query(query, tenderID, organizationID, authorUsername)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var review Review
		if err := rows.Scan(&review.ID, &review.BidID, &review.AuthorUsername, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
