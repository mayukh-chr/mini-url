package main

import (
	"net/http"
	"urlshortner/config"
	"urlshortner/database"
	"urlshortner/handlers"
	"urlshortner/middleware"
	"urlshortner/monitoring"
	"urlshortner/utils"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	cfg := config.Load()

	logger.WithFields(logrus.Fields{
		"environment": cfg.Environment,
		"port":        cfg.Port,
		"base_url":    cfg.BaseURL,
	}).Info("Starting URL Shortener server")

	// Initialize database
	database.InitDB(cfg.DatabaseURL)

	// Print URLs only in development
	if cfg.Environment == "development" {
		utils.PrintAllURLs()
	}

	// Create router with middleware
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter(100)) // 100 requests per second
	// Health check endpoint
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Monitoring endpoints
	r.HandleFunc("/metrics", monitoring.MetricsHandler).Methods("GET")
	r.HandleFunc("/metrics/prometheus", monitoring.PrometheusHandler).Methods("GET")

	// API routes
	r.HandleFunc("/shorten", handlers.ServeShortenPage).Methods("GET")
	r.HandleFunc("/shorten", handlers.CreateShortURL).Methods("POST")
	r.HandleFunc("/u/{code}", handlers.GetOriginalURL).Methods("GET")
	r.HandleFunc("/u/{code}", handlers.UpdateShortCode).Methods("PUT")
	r.HandleFunc("/u/{code}", handlers.DeleteShortURL).Methods("DELETE")
	r.HandleFunc("/stats/{code}", handlers.GetStats).Methods("GET")

	// Serve static files from frontend build
	if cfg.Environment == "production" {
		fs := http.FileServer(http.Dir("./frontend/build/"))
		r.PathPrefix("/").Handler(fs)
	}

	port := ":" + cfg.Port
	logger.WithField("port", port).Info("Server starting")

	if err := http.ListenAndServe(port, r); err != nil {
		logger.WithError(err).Fatal("Server failed to start")
	}
}
