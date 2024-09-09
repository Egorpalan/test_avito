package models

import (
	"database/sql"
	"time"
)

type Tender struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ServiceType     string    `json:"serviceType"`
	Status          string    `json:"status"`
	OrganizationID  int       `json:"organizationId"`
	CreatorUsername string    `json:"creatorUsername"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func GetAllTenders(db *sql.DB, serviceType string) ([]Tender, error) {
	var tenders []Tender
	var rows *sql.Rows
	var err error

	if serviceType == "" {
		rows, err = db.Query("SELECT id, name, description, service_type, status, organization_id, creator_username, created_at, updated_at FROM tenders")
	} else {
		rows, err = db.Query("SELECT id, name, description, service_type, status, organization_id, creator_username, created_at, updated_at FROM tenders WHERE service_type = $1", serviceType)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.CreatorUsername, &tender.CreatedAt, &tender.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func GetUserTenders(db *sql.DB, username string) ([]Tender, error) {
	var tenders []Tender
	query := `
        SELECT id, name, description, service_type, status, organization_id, creator_username, created_at, updated_at
        FROM tenders
        WHERE creator_username = $1
    `

	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tender Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationID, &tender.CreatorUsername, &tender.CreatedAt, &tender.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, tender)
	}

	return tenders, nil
}

func UpdateTender(db *sql.DB, id int, name, description string) (*Tender, error) {
	query := `
        UPDATE tenders
        SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
        RETURNING id, name, description, service_type, status, organization_id, creator_username, created_at, updated_at
    `

	var tender Tender
	err := db.QueryRow(query, name, description, id).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.CreatedAt,
		&tender.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tender, nil
}


func RollbackTender(db *sql.DB, id int, version int) (*Tender, error) {
	query := `
        SELECT id, name, description, service_type, status, organization_id, creator_username, created_at, updated_at
        FROM tenders
        WHERE id = $1 AND version = $2
    `
	var tender Tender
	err := db.QueryRow(query, id, version).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.CreatedAt,
		&tender.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tender, nil
}
