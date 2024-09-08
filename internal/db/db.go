package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq" // PostgreSQL driver
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
    return nil
}
