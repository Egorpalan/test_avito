package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := os.Getenv("POSTGRES_CONN")
	if connStr == "" {
		return fmt.Errorf("POSTGRES_CONN not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	DB = db
	if err := createTablesIfNotExists(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

func createTablesIfNotExists() error {
	var exists bool
	err := DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'tenders'
        )
    `).Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not check table existence: %v", err)
	}

	if !exists {
		_, err = DB.Exec(`
            CREATE TABLE tenders (
                id SERIAL PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                description TEXT,
                service_type VARCHAR(100),
                status VARCHAR(50),
                organization_id INT,
                creator_username VARCHAR(50),
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                version INT DEFAULT 1
            )
        `)
		if err != nil {
			return fmt.Errorf("could not create table: %v", err)
		}
	}

	err = DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'bids'
        )
    `).Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not check bids table existence: %v", err)
	}

	if !exists {
		_, err = DB.Exec(`
            CREATE TABLE bids (
                id SERIAL PRIMARY KEY,
                name VARCHAR(255) NOT NULL,
                description TEXT,
                status VARCHAR(50),
                tender_id INT REFERENCES tenders(id),
                organization_id INT,
                creator_username VARCHAR(50),
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        `)
		if err != nil {
			return fmt.Errorf("could not create bids table: %v", err)
		}
	}

	err = DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = 'reviews'
        )
    `).Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not check reviews table existence: %v", err)
	}

	if !exists {
		_, err = DB.Exec(`
            CREATE TABLE reviews (
                id SERIAL PRIMARY KEY,
                bid_id INT NOT NULL REFERENCES bids(id),
                author_username VARCHAR(50) NOT NULL,
                rating INT CHECK (rating >= 1 AND rating <= 5),
                comment TEXT,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        `)
		if err != nil {
			return fmt.Errorf("could not create reviews table: %v", err)
		}
	}

	return nil
}
