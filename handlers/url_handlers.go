package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"urlshortner/database"
	"urlshortner/models"
	"urlshortner/utils"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
)

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var u models.URL
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("Received CreateShortURL request: url=%s, short_code=%s", u.URL, u.ShortCode)
	if u.ShortCode == "" {
		u.ShortCode = utils.GenerateUniqueCode(6)
		log.Printf("Generated new short code: %s", u.ShortCode)
	} else {
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = ?)", u.ShortCode).Scan(&exists)
		if err != nil {
			log.Printf("Database error checking short code existence: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		if exists {
			log.Printf("Short code already exists: %s", u.ShortCode)
			http.Error(w, "short code already exists", http.StatusConflict)
			return
		}
	}

	stmt := `INSERT INTO urls (url, short_code) VALUES (?, ?)`
	_, err := database.DB.Exec(stmt, u.URL, u.ShortCode)
	if err != nil {
		log.Printf("Error inserting URL: %v", err)
		http.Error(w, "error inserting URL", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created short URL: %s -> %s", u.ShortCode, u.URL)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"short_code": u.ShortCode,
		"short_url":  "http://localhost:8080/u/" + u.ShortCode,
	}
	log.Printf("Responding with: %+v", resp)
	json.NewEncoder(w).Encode(resp)
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]

	row := database.DB.QueryRow(`SELECT url, access_count FROM urls WHERE short_code = ?`, shortCode)

	var url string
	var count int
	err := row.Scan(&url, &count)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		return
	}

	// Update access count
	_, _ = database.DB.Exec(`UPDATE urls SET access_count = access_count + 1 WHERE short_code = ?`, shortCode)

	http.Redirect(w, r, url, http.StatusFound)
}

func UpdateShortCode(w http.ResponseWriter, r *http.Request) {
	// Parse input JSON
	var payload struct {
		URL       string `json:"url"`
		ShortCode string `json:"short_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Prepare update statement: update short_code where url matches
	stmt := `UPDATE urls SET short_code = ?, updated_at = CURRENT_TIMESTAMP WHERE url = ?`

	res, err := database.DB.Exec(stmt, payload.ShortCode, payload.URL)
	if err != nil {
		// Check for UNIQUE constraint violation
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			http.Error(w, "Short code already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "URL not found or no changes made", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteShortURL(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]

	stmt := `DELETE FROM urls WHERE short_code = ?`
	_, err := database.DB.Exec(stmt, shortCode)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]

	row := database.DB.QueryRow(`SELECT access_count FROM urls WHERE short_code = ?`, shortCode)
	var count int
	if err := row.Scan(&count); err != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"access_count": count})
}

// ServeShortenPage serves the HTML form for shortening URLs
func ServeShortenPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/shorten.html")
}
