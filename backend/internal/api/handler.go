package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Unikyri/M-TurinChess/backend/internal/analysis"
	"github.com/Unikyri/M-TurinChess/backend/internal/db"
)

// Handler groups the HTTP handlers that depend on the Analyzer and Storage.
type Handler struct {
	analyzer *analysis.Analyzer
	storage  db.Storage
}

// NewHandler creates a Handler with the given dependencies.
func NewHandler(analyzer *analysis.Analyzer, storage db.Storage) *Handler {
	return &Handler{analyzer: analyzer, storage: storage}
}

// HandleAnalyze handles POST /api/analyze.
// Accepts a multipart form with a pgn_file and optional fields.
func (h *Handler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10 MB).
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, "cannot parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// --- Required fields ---

	pgnText := r.FormValue("pgn_text")
	if pgnText == "" {
		jsonError(w, "pgn_text is required", http.StatusBadRequest)
		return
	}
	pgnBytes := []byte(pgnText)

	playerColor := r.FormValue("player_color")
	if playerColor != "white" && playerColor != "black" {
		jsonError(w, "player_color must be 'white' or 'black'", http.StatusBadRequest)
		return
	}

	// --- Optional integer fields ---
	eloWhite := parseIntField(r, "elo_white", 0)
	eloBlack := parseIntField(r, "elo_black", 0)
	threshold := parseIntField(r, "threshold", 6)
	depth := parseIntField(r, "depth", 18) // Restored to 18 for final analysis precision

	req := analysis.AnalysisRequest{
		PGNBytes:    pgnBytes,
		PlayerColor: playerColor,
		EloWhite:    eloWhite,
		EloBlack:    eloBlack,
		Threshold:   threshold,
		Depth:       depth,
	}

	result, err := h.analyzer.Analyze(req)
	if err != nil {
		log.Printf("[analyze] error: %v", err)
		// PGN parsing errors are client mistakes → 400; everything else → 500.
		if strings.Contains(err.Error(), "parse pgn") {
			jsonError(w, "invalid pgn file: "+err.Error(), http.StatusBadRequest)
		} else {
			jsonError(w, "analysis failed: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("[analyze] encode response: %v", err)
	}
}

// HandleHistory handles GET /api/history.
// Returns all analysis summaries stored in the JSON database.
func (h *Handler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	records, err := h.storage.ReadAll()
	if err != nil {
		log.Printf("[history] read error: %v", err)
		jsonError(w, "cannot read history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(records); err != nil {
		log.Printf("[history] encode response: %v", err)
	}
}

// --- Helpers ---

// jsonError writes a JSON-encoded error response.
func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		log.Printf("[jsonError] encode: %v", err)
	}
}

// parseIntField reads a form value as an integer, returning defaultVal on failure.
func parseIntField(r *http.Request, field string, defaultVal int) int {
	v := r.FormValue(field)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return defaultVal
	}
	return n
}
