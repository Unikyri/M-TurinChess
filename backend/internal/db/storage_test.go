package db

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestJSONStorage_ReadAll_NonExistentFile(t *testing.T) {
	s := NewJSONStorage("/tmp/this_file_does_not_exist_mturin_xyz123.json")
	records, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on non-existent file: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty slice, got %d records", len(records))
	}
}

func TestJSONStorage_AppendAndReadAll(t *testing.T) {
	f, err := os.CreateTemp("", "history_test_*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	s := NewJSONStorage(f.Name())

	// Initially empty (file exists but is empty).
	records, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on empty file: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}

	// Append one record.
	rec := HistoryRecord{
		ID:                 "abc123",
		AnalyzedAt:         time.Now().UTC(),
		Verdict:            "HUMAN_PLAYER",
		SuspicionCount:     2,
		Threshold:          6,
		TotalMovesAnalyzed: 20,
		PlayerColor:        "white",
		Elo:                1500,
	}
	if err := s.Append(rec); err != nil {
		t.Fatalf("Append: %v", err)
	}

	// Read back and verify.
	records, err = s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll after Append: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	got := records[0]
	if got.ID != "abc123" {
		t.Errorf("ID: got %q, want abc123", got.ID)
	}
	if got.Verdict != "HUMAN_PLAYER" {
		t.Errorf("Verdict: got %q, want HUMAN_PLAYER", got.Verdict)
	}
	if got.Elo != 1500 {
		t.Errorf("Elo: got %d, want 1500", got.Elo)
	}
}

func TestJSONStorage_AppendMultiple(t *testing.T) {
	f, err := os.CreateTemp("", "history_test_*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Close()
	defer os.Remove(f.Name())

	s := NewJSONStorage(f.Name())

	for i := 0; i < 5; i++ {
		rec := HistoryRecord{
			ID:      fmt.Sprintf("id%d", i),
			Verdict: "HUMAN_PLAYER",
		}
		if err := s.Append(rec); err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
	}

	records, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(records) != 5 {
		t.Errorf("expected 5 records, got %d", len(records))
	}
	for i, r := range records {
		want := fmt.Sprintf("id%d", i)
		if r.ID != want {
			t.Errorf("records[%d].ID: got %q, want %q", i, r.ID, want)
		}
	}
}
