package handlers

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/gorilla/mux"
	"urlshortner/database"
	"urlshortner/models"
)


func CreateShortURL(w http.ResponseWriter, r *http.Request){
	var u models.URL
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil{
		http.Error(w, "invalid input", http.StatusBadRequest)
		return 
	}

	stmt := `INSERT INTO urls (url, short_code) VALUES (?, ?)`

	_, err := database.DB.Exec(stmt, u.URL, u.ShortCode)

	if err != nil {
		http.Error(w, "error inserting URL", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetOriginalURL(w http.ResponseWriter, r* http.Request){
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

func UpdateShortURL(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]
	var u models.URL
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	stmt := `UPDATE urls SET url = ?, updated_at = CURRENT_TIMESTAMP WHERE short_code = ?`
	_, err := database.DB.Exec(stmt, u.URL, shortCode)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
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