// d:\Code\repos\golang-tutorial\url-shortner\utils\utils.go
package utils

import (
	"math/rand"
	"time"
	"urlshortner/database"
	"fmt"
	"log"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Exported version for use in handlers
func GenerateUniqueCode(length int) string {
	for {
		code := generateRandomCode(length)
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = ?)", code).Scan(&exists)
		if err != nil {
			// Log error and fallback
			continue
		}
		if !exists {
			return code
		}
	}
}

//testing 

func PrintAllURLs() {
	database.InitDB("urlshortener.db")

	rows, err := database.DB.Query("SELECT id, url, short_code, access_count FROM urls")
	if err != nil {
		log.Fatalf("Error querying database: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("All URLs in database:")
	for rows.Next() {
		var id int
		var url, shortCode string
		var count int

		if err := rows.Scan(&id, &url, &shortCode, &count); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		fmt.Printf("ID: %d | URL: %s | ShortCode: %s | AccessCount: %d\n", id, url, shortCode, count)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}