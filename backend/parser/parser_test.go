package parser

import (
	"testing"
)

func TestParser_Valid(t *testing.T) {
	validFENs := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-",
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR|b|KQkq|e3",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|-|f6",
		"8/8/8/8/8/8/8/8|b|-|-",
	}

	for _, fen := range validFENs {
		t.Run(fen, func(t *testing.T) {
			lexer := NewLexer(fen)
			parser := NewParser(lexer)
			err := parser.Parse()
			if err != nil {
				t.Errorf("expected no error for valid FEN %q, got: %v", fen, err)
			}
		})
	}
}

func TestParser_Invalid(t *testing.T) {
	invalidFENs := []struct {
		fen string
		msg string
	}{
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|x|KQkq|-",
			msg: "expected w or b",
		},
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkqQ|-",
			msg: "duplicate castling",
		},
		{
			fen: "rnbqkbnr/pppppppp/7/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-",
			msg: "row sum is 7",
		},
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP|w|KQkq|-",
			msg: "expected exactly 8 rows",
		},
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|e9",
			msg: "invalid en passant",
		},
		{
			fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq",
			msg: "missing fields",
		},
	}

	for _, tc := range invalidFENs {
		t.Run(tc.fen, func(t *testing.T) {
			lexer := NewLexer(tc.fen)
			parser := NewParser(lexer)
			err := parser.Parse()
			if err == nil {
				t.Errorf("expected error for invalid FEN %q, got none", tc.fen)
			}
		})
	}
}
