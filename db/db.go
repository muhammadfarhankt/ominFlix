package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DB struct to hold the database connection
type DB struct {
	DB *sql.DB
}

// NewDB creates a new DB instance, loading env vars and connecting to the database
func NewDB() (*DB, error) {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")

	dbPass := os.Getenv("DB_PASS")

	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &DB{DB: db}, nil
}

// CreateTable creates the 'blocks' table if it doesn't exist
func (d *DB) CreateTable() error {
	_, err := d.DB.Exec(`CREATE TABLE IF NOT EXISTS blocks (
        block_height BIGINT PRIMARY KEY,
		block_id TEXT,
        proposer_address TEXT,
        num_transactions INT,
        details JSONB,
        created_at TIMESTAMP WITH TIME ZONE,
        updated_at TIMESTAMP WITH TIME ZONE,
        deleted_at TIMESTAMP WITH TIME ZONE
      )`)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	// Create index on block_height
	_, err = d.DB.Exec(`CREATE INDEX IF NOT EXISTS blocks_height_idx ON blocks (block_height)`)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *DB) Close() {
	d.DB.Close()
}
