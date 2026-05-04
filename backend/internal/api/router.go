package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// NewRouter creates the HTTP mux, registers all routes, and wraps with CORS.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// Health probe — no auth, no CORS needed, used by container health checks.
	mux.HandleFunc("/health", handleHealth)

	// API routes.
	mux.HandleFunc("/api/analyze", h.HandleAnalyze)
	mux.HandleFunc("/api/history", h.HandleHistory)

	return corsMiddleware(mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("[health] encode: %v", err)
	}
}
