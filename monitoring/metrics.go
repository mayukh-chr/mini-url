package monitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"urlshortner/database"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
}

type MetricsResponse struct {
	System   SystemMetrics   `json:"system"`
	Database DatabaseMetrics `json:"database"`
	App      AppMetrics      `json:"app"`
}

type SystemMetrics struct {
	Uptime        string `json:"uptime"`
	GoVersion     string `json:"go_version"`
	NumGoroutines int    `json:"num_goroutines"`
	MemoryAlloc   uint64 `json:"memory_alloc_mb"`
	MemoryTotal   uint64 `json:"memory_total_mb"`
	MemorySys     uint64 `json:"memory_sys_mb"`
	NumGC         uint32 `json:"num_gc"`
}

type DatabaseMetrics struct {
	Connected  bool `json:"connected"`
	OpenConns  int  `json:"open_connections"`
	InUseConns int  `json:"in_use_connections"`
	IdleConns  int  `json:"idle_connections"`
}

type AppMetrics struct {
	Version     string    `json:"version"`
	Environment string    `json:"environment"`
	Timestamp   time.Time `json:"timestamp"`
}

var startTime = time.Now()

// MetricsHandler provides detailed application metrics
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := database.DB.Stats()

	metrics := MetricsResponse{
		System: SystemMetrics{
			Uptime:        time.Since(startTime).String(),
			GoVersion:     runtime.Version(),
			NumGoroutines: runtime.NumGoroutine(),
			MemoryAlloc:   bToMb(m.Alloc),
			MemoryTotal:   bToMb(m.TotalAlloc),
			MemorySys:     bToMb(m.Sys),
			NumGC:         m.NumGC,
		},
		Database: DatabaseMetrics{
			Connected:  database.DB.Ping() == nil,
			OpenConns:  stats.OpenConnections,
			InUseConns: stats.InUse,
			IdleConns:  stats.Idle,
		},
		App: AppMetrics{
			Version:     "1.0.0",
			Environment: "production", // This should come from config
			Timestamp:   time.Now().UTC(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// PrometheusHandler provides metrics in Prometheus format
func PrometheusHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := database.DB.Stats()
	uptime := time.Since(startTime).Seconds()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	metrics := `# HELP urlshortener_uptime_seconds Total uptime in seconds
# TYPE urlshortener_uptime_seconds counter
urlshortener_uptime_seconds %f

# HELP urlshortener_memory_alloc_bytes Currently allocated memory in bytes
# TYPE urlshortener_memory_alloc_bytes gauge
urlshortener_memory_alloc_bytes %d

# HELP urlshortener_goroutines Number of goroutines
# TYPE urlshortener_goroutines gauge
urlshortener_goroutines %d

# HELP urlshortener_db_connections Database connection pool stats
# TYPE urlshortener_db_connections gauge
urlshortener_db_connections{state="open"} %d
urlshortener_db_connections{state="in_use"} %d
urlshortener_db_connections{state="idle"} %d

# HELP urlshortener_gc_count Total number of garbage collections
# TYPE urlshortener_gc_count counter
urlshortener_gc_count %d
`

	w.Write([]byte(fmt.Sprintf(metrics,
		uptime,
		m.Alloc,
		runtime.NumGoroutine(),
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		m.NumGC,
	)))
}
