package database

import (
	"database/sql"
	"log"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// creates db if it doesn't exist.
func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	
	fmt.Println("db opened")

	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		access_count INTEGER NOT NULL DEFAULT 0
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
