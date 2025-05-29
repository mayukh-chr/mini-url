package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB
var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
}

// InitDB initializes database connection with PostgreSQL for production or SQLite for development
func InitDB(databaseURL string) {
	var err error

	if databaseURL == "" {
		// Development mode - use SQLite
		logger.Info("Initializing SQLite database for development")
		DB, err = sql.Open("sqlite3", "./urlshortener.db")
		if err != nil {
			log.Fatalf("Failed to open SQLite database: %v", err)
		}
		createSQLiteTables()
	} else {
		// Production mode - use PostgreSQL
		logger.Info("Initializing PostgreSQL database for production")
		DB, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Fatalf("Failed to open PostgreSQL database: %v", err)
		}

		// Configure connection pool
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(25)
		DB.SetConnMaxLifetime(5 * time.Minute)

		createPostgresTables()
	}

	// Test connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	logger.Info("Database connection established successfully")
}

func createSQLiteTables() {
	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		access_count INTEGER NOT NULL DEFAULT 0
	);`

	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create SQLite table: %v", err)
	}
}

func createPostgresTables() {
	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		short_code VARCHAR(50) UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		access_count INTEGER NOT NULL DEFAULT 0
	);
	
	CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
	CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);`

	_, err := DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL tables: %v", err)
	}
}

// HealthCheck verifies database connectivity
func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return DB.PingContext(ctx)
}
