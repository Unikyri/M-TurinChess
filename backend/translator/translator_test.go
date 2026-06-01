package translator

import (
	"testing"
)

func TestSimplifyFEN(t *testing.T) {
	standard := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-"

	simplified, err := SimplifyFEN(standard)
	if err != nil {
		t.Fatalf("SimplifyFEN returned error: %v", err)
	}

	if simplified != expected {
		t.Errorf("expected simplified FEN %q, got %q", expected, simplified)
	}
}

func TestTranslatePGN(t *testing.T) {
	pgn := "1. e4 e5 2. Nf3 Nc6"
	fens, err := TranslatePGN(pgn)
	if err != nil {
		t.Fatalf("TranslatePGN returned error: %v", err)
	}

	// 2 moves = 3 positions (initial, after e4, after e5, after Nf3, after Nc6? Wait, e4 e5 Nf3 Nc6 is 4 moves!)
	// So 4 moves = 5 positions
	if len(fens) != 5 {
		t.Errorf("expected 5 positions for 4 moves, got %d", len(fens))
	}

	// Initial FEN should be starting FEN
	expectedStart := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-"
	if fens[0] != expectedStart {
		t.Errorf("expected starting FEN %q, got %q", expectedStart, fens[0])
	}
}

func TestTranslatePGNCastling(t *testing.T) {
	pgn := "1. e4 e5 2. Nf3 Nc6 3. Bc4 Bc5 4. 0-0 Nf6 5. d3 0-0"
	fens, err := TranslatePGN(pgn)
	if err != nil {
		t.Fatalf("TranslatePGN with 0-0 castling failed: %v", err)
	}
	if len(fens) == 0 {
		t.Errorf("expected some FEN positions, got 0")
	}
}
