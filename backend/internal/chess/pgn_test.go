package chess

import (
	"testing"
)

// --- Test fixtures ---

const pgnWithElo = `[Event "F/S Return Match"]
[Site "Belgrade, Serbia JUG"]
[Date "1992.11.04"]
[White "Fischer, Robert J."]
[Black "Spassky, Boris V."]
[WhiteElo "2785"]
[BlackElo "2740"]
[Result "1/2-1/2"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 1/2-1/2`

const pgnNoElo = `[White "Player One"]
[Black "Player Two"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 1-0`

const pgnWithCommentsAndVariations = `[White "A"]
[Black "B"]
[Result "*"]

1. e4 {This is a Ruy Lopez comment} e5 (1... c5 {Sicilian} 2. Nf3 d6) 2. Nf3 Nc6 *`

const pgnNestedVariations = `[White "A"]
[Black "B"]
[Result "*"]

1. e4 e5 (1... c5 (1... e6 2. d4) 2. Nf3) 2. Nf3 *`

const pgnWithNAGs = `[White "A"]
[Black "B"]
[Result "*"]

1. e4 $1 e5 $2 2. Nf3 $3 Nc6 *`

const pgnEmptyMovetext = `[White "A"]
[Black "B"]
[Result "*"]

*`

// --- Tests ---

func TestParsePGN_FullMetadata(t *testing.T) {
	g, err := ParsePGN([]byte(pgnWithElo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"White", g.White, "Fischer, Robert J."},
		{"Black", g.Black, "Spassky, Boris V."},
		{"WhiteElo", g.WhiteElo, 2785},
		{"BlackElo", g.BlackElo, 2740},
		{"Result", g.Result, "1/2-1/2"},
		{"Event", g.Event, "F/S Return Match"},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestParsePGN_NoElo(t *testing.T) {
	g, err := ParsePGN([]byte(pgnNoElo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.WhiteElo != 0 {
		t.Errorf("WhiteElo: got %d, want 0", g.WhiteElo)
	}
	if g.BlackElo != 0 {
		t.Errorf("BlackElo: got %d, want 0", g.BlackElo)
	}
}

func TestParsePGN_CommentsStripped(t *testing.T) {
	g, err := ParsePGN([]byte(pgnWithCommentsAndVariations))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"e4", "e5", "Nf3", "Nc6"}
	assertMoves(t, g.Moves, want)
}

func TestParsePGN_NestedVariationsStripped(t *testing.T) {
	g, err := ParsePGN([]byte(pgnNestedVariations))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"e4", "e5", "Nf3"}
	assertMoves(t, g.Moves, want)
}

func TestParsePGN_NAGsStripped(t *testing.T) {
	g, err := ParsePGN([]byte(pgnWithNAGs))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"e4", "e5", "Nf3", "Nc6"}
	assertMoves(t, g.Moves, want)
}

func TestParsePGN_EmptyFile(t *testing.T) {
	_, err := ParsePGN([]byte{})
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
}

func TestParsePGN_EmptyMovetext(t *testing.T) {
	g, err := ParsePGN([]byte(pgnEmptyMovetext))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Moves) != 0 {
		t.Errorf("Moves: got %v, want empty slice", g.Moves)
	}
}

func TestParsePGN_ResultNotInMoves(t *testing.T) {
	g, err := ParsePGN([]byte(pgnNoElo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, m := range g.Moves {
		if gameTerminators[m] {
			t.Errorf("result token %q found in Moves slice", m)
		}
	}
}

func TestParsePGN_MoveNumbersNotInMoves(t *testing.T) {
	g, err := ParsePGN([]byte(pgnWithElo))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, m := range g.Moves {
		if isMoveNumber(m) {
			t.Errorf("move number token %q found in Moves slice", m)
		}
	}
}

func TestParsePGN_CastlingPreserved(t *testing.T) {
	pgn := `[White "A"]
[Black "B"]
[Result "*"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 *`

	g, err := ParsePGN([]byte(pgn))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, m := range g.Moves {
		if m == "O-O" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("castling O-O not found in Moves: %v", g.Moves)
	}
}

// --- Helpers ---

func assertMoves(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("Moves length: got %d (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Moves[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}
