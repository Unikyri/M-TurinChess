package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// HistoryRecord is a compact record stored in history.json.
// It omits move details and traces to keep the file size manageable.
type HistoryRecord struct {
	ID                 string    `json:"id"`
	AnalyzedAt         time.Time `json:"analyzed_at"`
	Verdict            string    `json:"verdict"`
	SuspicionCount     int       `json:"suspicion_count"`
	Threshold          int       `json:"threshold"`
	TotalMovesAnalyzed int       `json:"total_moves_analyzed"`
	PlayerColor        string    `json:"player_color"`
	Elo                int       `json:"elo"`
}

// Storage is the persistence interface for analysis records.
type Storage interface {
	Append(record HistoryRecord) error
	ReadAll() ([]HistoryRecord, error)
}

// JSONStorage implements Storage using a single JSON file.
type JSONStorage struct {
	path string
	mu   sync.Mutex
}

// NewJSONStorage creates a JSONStorage backed by the file at path.
// The file is created on first write; reading a missing file returns an empty slice.
func NewJSONStorage(path string) *JSONStorage {
	return &JSONStorage{path: path}
}

// ReadAll returns all history records. Returns an empty slice if the file does not exist.
func (s *JSONStorage) ReadAll() ([]HistoryRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readLocked()
}

// Append appends a single record to the history file.
func (s *JSONStorage) Append(record HistoryRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	records, err := s.readLocked()
	if err != nil {
		return fmt.Errorf("read before append: %w", err)
	}
	records = append(records, record)

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal records: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("write history file: %w", err)
	}
	return nil
}

// readLocked reads the JSON file. Must be called with s.mu held.
func (s *JSONStorage) readLocked() ([]HistoryRecord, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []HistoryRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history file: %w", err)
	}
	if len(data) == 0 {
		return []HistoryRecord{}, nil
	}
	var records []HistoryRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("unmarshal history: %w", err)
	}
	return records, nil
}
