package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// HealthHandler holds dependencies needed for readiness checks.
type HealthHandler struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

// Healthz is a simple liveness probe: it only indicates the process is running.
func (h *HealthHandler) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// Readyz checks whether critical dependencies are reachable.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	type status struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}

	// Check PostgreSQL if configured.
	if h.DB != nil {
		if err := h.DB.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(status{Status: "not ready", Error: "database unavailable"})
			return
		}
	}

	// Check Redis if configured.
	if h.Redis != nil {
		if err := h.Redis.Ping(r.Context()).Err(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(status{Status: "not ready", Error: "redis unavailable"})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

