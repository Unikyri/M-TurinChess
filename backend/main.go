package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"backend/parser"
	"backend/translator"
	"backend/turing"

	"github.com/corentings/chess/v2"
)

type AnalyzeRequest struct {
	Pgn string `json:"pgn"`
}

type SimulateRequest struct {
	History []string `json:"history"`
	Current string   `json:"current"`
}

type AnalyzeResponse struct {
	Moves               []string                    `json:"moves"`
	SimplifiedFens      []string                    `json:"simplifiedFens"`
	Simulations         []*turing.SimulationResult `json:"simulations"`
	RepetitionDetected  bool                        `json:"repetitionDetected"`
	RepetitionMoveIndex int                         `json:"repetitionMoveIndex"` // -1 if not detected, else 0-indexed move index
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/analyze", handleAnalyze)
	mux.HandleFunc("/api/simulate", handleSimulate)

	handler := enableCORS(mux)

	port := ":8080"
	log.Printf("Backend server listening on http://localhost%s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Pgn) == "" {
		writeError(w, "PGN string cannot be empty", http.StatusBadRequest)
		return
	}

	// 1. Translate PGN to simplified FENs
	simplifiedFens, err := translator.TranslatePGN(req.Pgn)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to translate PGN: %v", err), http.StatusBadRequest)
		return
	}

	// 2. Validate all FENs using Lexer and Parser
	for i, fen := range simplifiedFens {
		lexer := parser.NewLexer(fen)
		p := parser.NewParser(lexer)
		if err := p.Parse(); err != nil {
			writeError(w, fmt.Sprintf("Validation error at position %d: %v", i, err), http.StatusBadRequest)
			return
		}
	}

	// 3. Extract SAN moves list from PGN
	pgnReader := strings.NewReader(req.Pgn)
	pgnOpt, err := chess.PGN(pgnReader)
	if err != nil {
		writeError(w, fmt.Sprintf("Failed to parse PGN for moves: %v", err), http.StatusBadRequest)
		return
	}
	game := chess.NewGame(pgnOpt)
	positions := game.Positions()
	moves := game.Moves()

	var notation chess.Notation = chess.AlgebraicNotation{}
	sanMoves := make([]string, len(moves))
	for i, m := range moves {
		sanMoves[i] = notation.Encode(positions[i], m)
	}

	// 4. Run Turing Machine for each move
	simulations := make([]*turing.SimulationResult, len(moves))
	repetitionDetected := false
	repetitionMoveIndex := -1

	for i := 0; i < len(moves); i++ {
		// History up to i is simplifiedFens[0 : i+1]
		// Current FEN is simplifiedFens[i+1]
		history := simplifiedFens[0 : i+1]
		current := simplifiedFens[i+1]

		simResult := turing.Simulate(history, current)

		if simResult.Accepted && !repetitionDetected {
			repetitionDetected = true
			repetitionMoveIndex = i
		}

		// Omit steps in the analyze response to prevent OOM
		simResult.Steps = nil
		simulations[i] = simResult
	}

	resp := AnalyzeResponse{
		Moves:               sanMoves,
		SimplifiedFens:      simplifiedFens,
		Simulations:         simulations,
		RepetitionDetected:  repetitionDetected,
		RepetitionMoveIndex: repetitionMoveIndex,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SimulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Current == "" {
		writeError(w, "Current FEN string cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate History FENs
	for i, fen := range req.History {
		lexer := parser.NewLexer(fen)
		p := parser.NewParser(lexer)
		if err := p.Parse(); err != nil {
			writeError(w, fmt.Sprintf("Validation error in history at position %d: %v", i, err), http.StatusBadRequest)
			return
		}
	}

	// Validate Current FEN
	lexer := parser.NewLexer(req.Current)
	p := parser.NewParser(lexer)
	if err := p.Parse(); err != nil {
		writeError(w, fmt.Sprintf("Validation error in current FEN: %v", err), http.StatusBadRequest)
		return
	}

	simResult := turing.Simulate(req.History, req.Current)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(simResult)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
