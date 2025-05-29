package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"regexp"
	"strings"
	"time"

	"urlshortner/database"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateUniqueCode(length int) string {
	maxAttempts := 10
	for attempt := 0; attempt < maxAttempts; attempt++ {
		code := generateRandomCode(length)
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = ?)", code).Scan(&exists)
		if err != nil {
			log.Printf("Error checking code uniqueness: %v", err)
			continue
		}
		if !exists {
			return code
		}
	}
	// If we can't find a unique code, increase length
	return GenerateUniqueCode(length + 1)
}

// IsValidURL validates if the provided string is a valid URL
func IsValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	// Basic regex check
	if !urlRegex.MatchString(urlStr) {
		return false
	}

	// Parse URL to ensure it's well-formed
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Ensure scheme and host are present
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

// IsValidShortCode validates short code format
func IsValidShortCode(code string) bool {
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	// Only allow alphanumeric characters
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9')) {
			return false
		}
	}

	return true
}

// SanitizeURL cleans and normalizes URL
func SanitizeURL(urlStr string) string {
	urlStr = strings.TrimSpace(urlStr)

	// Add http:// if no scheme is provided
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}

	return urlStr
}

// PrintAllURLs - testing function
func PrintAllURLs() {
	rows, err := database.DB.Query("SELECT id, url, short_code, access_count FROM urls")
	if err != nil {
		log.Printf("Error querying database: %v\n", err)
		return
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
