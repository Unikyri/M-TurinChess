package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration_Acceptance(t *testing.T) {
	handler := http.HandlerFunc(handleAnalyze)

	// Test Case 1: Acceptance (Triple Repetition)
	acceptPgn := "1. Nf3 Nf6 2. Ng1 Ng8 3. Nf3 Nf6 4. Ng1 Ng8 5. Nf3 Nf6 6. Ng1 Ng8"
	reqBody, _ := json.Marshal(AnalyzeRequest{Pgn: acceptPgn})
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Test Case 1 returned status %d, expected %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp AnalyzeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response for Test Case 1: %v", err)
	}

	if !resp.RepetitionDetected {
		t.Fatal("Test Case 1 did not detect triple repetition!")
	}
	t.Logf("Test Case 1 (Acceptance) passed. Repetition detected at move index %d (move: %s)", resp.RepetitionMoveIndex, resp.Moves[resp.RepetitionMoveIndex])
}

func TestIntegration_Rejection(t *testing.T) {
	handler := http.HandlerFunc(handleAnalyze)

	// Test Case 2: Rejection (No Repetition)
	rejectPgn := "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6"
	reqBody, _ := json.Marshal(AnalyzeRequest{Pgn: rejectPgn})
	req := httptest.NewRequest(http.MethodPost, "/api/analyze", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Test Case 2 returned status %d, expected %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp AnalyzeResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response for Test Case 2: %v", err)
	}

	if resp.RepetitionDetected {
		t.Fatal("Test Case 2 incorrectly detected triple repetition!")
	}
	t.Log("Test Case 2 (Rejection) passed. No repetition detected.")
}

func TestIntegration_Simulate(t *testing.T) {
	handler := http.HandlerFunc(handleSimulate)

	// Simulate a simple transition
	history := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-",
		"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R|b|KQkq|-",
	}
	current := "rnbqkb1r/pppppppp/5n2/8/8/5N2/PPPPPPPP/RNBQKB1R|w|KQkq|-"

	reqBody, _ := json.Marshal(SimulateRequest{History: history, Current: current})
	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Simulate endpoint returned status %d, expected %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var resp struct {
		Steps      []interface{} `json:"steps"`
		FinalState string        `json:"finalState"`
		Accepted   bool          `json:"accepted"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode simulate response: %v", err)
	}

	if len(resp.Steps) == 0 {
		t.Fatal("Simulate response did not return any steps!")
	}

	t.Logf("Simulate integration test passed. Total steps simulated: %d, Final state: %s, Accepted: %v", len(resp.Steps), resp.FinalState, resp.Accepted)
}

