package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"urlshortner/config"
	"urlshortner/database"
	"urlshortner/models"
	"urlshortner/utils"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var cfg *config.Config

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	cfg = config.Load()
}

func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var u models.URL
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logger.WithError(err).Warn("Invalid JSON input")
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Sanitize and validate URL
	u.URL = utils.SanitizeURL(u.URL)
	if !utils.IsValidURL(u.URL) {
		logger.WithField("url", u.URL).Warn("Invalid URL provided")
		http.Error(w, "invalid URL format", http.StatusBadRequest)
		return
	}

	logger.WithFields(logrus.Fields{
		"url":        u.URL,
		"short_code": u.ShortCode,
	}).Info("Received CreateShortURL request")

	if u.ShortCode == "" {
		u.ShortCode = utils.GenerateUniqueCode(6)
		logger.WithField("short_code", u.ShortCode).Info("Generated new short code")
	} else {
		// Validate custom short code
		if !utils.IsValidShortCode(u.ShortCode) {
			logger.WithField("short_code", u.ShortCode).Warn("Invalid short code format")
			http.Error(w, "invalid short code format", http.StatusBadRequest)
			return
		}

		var exists bool
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := database.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)", u.ShortCode).Scan(&exists)
		if err != nil {
			logger.WithError(err).Error("Database error checking short code existence")
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		if exists {
			logger.WithField("short_code", u.ShortCode).Warn("Short code already exists")
			http.Error(w, "short code already exists", http.StatusConflict)
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `INSERT INTO urls (url, short_code) VALUES ($1, $2)`
	_, err := database.DB.ExecContext(ctx, stmt, u.URL, u.ShortCode)
	if err != nil {
		// Handle PostgreSQL constraint violations
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "short code already exists", http.StatusConflict)
			return
		}
		// Handle SQLite constraint violations
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			http.Error(w, "short code already exists", http.StatusConflict)
			return
		}

		logger.WithError(err).Error("Error inserting URL")
		http.Error(w, "error inserting URL", http.StatusInternalServerError)
		return
	}

	logger.WithFields(logrus.Fields{
		"short_code": u.ShortCode,
		"url":        u.URL,
	}).Info("Successfully created short URL")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]string{
		"short_code": u.ShortCode,
		"short_url":  cfg.BaseURL + "/u/" + u.ShortCode,
	}
	json.NewEncoder(w).Encode(resp)
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := database.DB.QueryRowContext(ctx, `SELECT url, access_count FROM urls WHERE short_code = $1`, shortCode)

	var url string
	var count int
	err := row.Scan(&url, &count)

	if err == sql.ErrNoRows {
		logger.WithField("short_code", shortCode).Warn("Short code not found")
		http.NotFound(w, r)
		return
	} else if err != nil {
		logger.WithError(err).Error("Error fetching URL")
		http.Error(w, "Error fetching URL", http.StatusInternalServerError)
		return
	}

	// Update access count asynchronously
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := database.DB.ExecContext(ctx, `UPDATE urls SET access_count = access_count + 1, updated_at = CURRENT_TIMESTAMP WHERE short_code = $1`, shortCode)
		if err != nil {
			logger.WithError(err).Error("Failed to update access count")
		}
	}()

	logger.WithFields(logrus.Fields{
		"short_code":   shortCode,
		"redirect_url": url,
		"access_count": count + 1,
	}).Info("Redirecting user")

	http.Redirect(w, r, url, http.StatusFound)
}

func UpdateShortCode(w http.ResponseWriter, r *http.Request) {
	// Parse input JSON
	var payload struct {
		URL       string `json:"url"`
		ShortCode string `json:"short_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.WithError(err).Warn("Invalid JSON input for update")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if !utils.IsValidURL(payload.URL) {
		logger.WithField("url", payload.URL).Warn("Invalid URL in update request")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	if !utils.IsValidShortCode(payload.ShortCode) {
		logger.WithField("short_code", payload.ShortCode).Warn("Invalid short code in update request")
		http.Error(w, "Invalid short code format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare update statement: update short_code where url matches
	stmt := `UPDATE urls SET short_code = $1, updated_at = CURRENT_TIMESTAMP WHERE url = $2`

	res, err := database.DB.ExecContext(ctx, stmt, payload.ShortCode, payload.URL)
	if err != nil {
		// Check for UNIQUE constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			http.Error(w, "Short code already exists", http.StatusConflict)
			return
		}
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			http.Error(w, "Short code already exists", http.StatusConflict)
			return
		}

		logger.WithError(err).Error("Database error during update")
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logger.WithFields(logrus.Fields{
			"url":        payload.URL,
			"short_code": payload.ShortCode,
		}).Warn("URL not found for update")
		http.Error(w, "URL not found or no changes made", http.StatusNotFound)
		return
	}

	logger.WithFields(logrus.Fields{
		"url":        payload.URL,
		"short_code": payload.ShortCode,
	}).Info("Successfully updated short code")

	w.WriteHeader(http.StatusOK)
}

func DeleteShortURL(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `DELETE FROM urls WHERE short_code = $1`
	res, err := database.DB.ExecContext(ctx, stmt, shortCode)
	if err != nil {
		logger.WithError(err).Error("Database error during delete")
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		logger.WithField("short_code", shortCode).Warn("Short code not found for deletion")
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	logger.WithField("short_code", shortCode).Info("Successfully deleted short URL")
	w.WriteHeader(http.StatusOK)
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["code"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := database.DB.QueryRowContext(ctx, `SELECT access_count FROM urls WHERE short_code = $1`, shortCode)
	var count int
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("short_code", shortCode).Warn("Short code not found for stats")
			http.NotFound(w, r)
			return
		}
		logger.WithError(err).Error("Database error fetching stats")
		http.Error(w, "Error fetching stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"access_count": count})
}

// HealthCheck endpoint for monitoring
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connectivity
	if err := database.HealthCheck(); err != nil {
		logger.WithError(err).Error("Health check failed - database unreachable")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "database unreachable",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "healthy",
		"timestamp":   time.Now().UTC(),
		"version":     "1.0.0",
		"environment": cfg.Environment,
	})
}

// ServeShortenPage serves the HTML form for shortening URLs
func ServeShortenPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/shorten.html")
}
