package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// PostgreSQL connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "currencies"),
	)

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create table if not exists
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS currencies (
			id SERIAL PRIMARY KEY,
			currency_code VARCHAR(10) NOT NULL,
			buy_rate DECIMAL(10,4) NOT NULL,
			sell_rate DECIMAL(10,4) NOT NULL,
			timestamp TIMESTAMP NOT NULL
		)
	`); err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	// Create index for faster queries
	if _, err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_currency_timestamp ON currencies (currency_code, timestamp)
	`); err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	log.Println("PostgreSQL database connection established")
	return db, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}