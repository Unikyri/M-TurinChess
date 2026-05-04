package chess

import "testing"

func TestBoard_InitialPosition(t *testing.T) {
	b := NewBoard()
	// White king at e1 (square 4)
	if b.squares[4].Type != King || b.squares[4].Color != ColorWhite {
		t.Errorf("expected white king at e1")
	}
	// Black queen at d8 (square 59)
	if b.squares[59].Type != Queen || b.squares[59].Color != ColorBlack {
		t.Errorf("expected black queen at d8")
	}
}

func TestBoard_SANToUCI_PawnE4(t *testing.T) {
	b := NewBoard()
	uci, err := b.SANToUCI("e4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uci != "e2e4" {
		t.Errorf("got %q, want %q", uci, "e2e4")
	}
}

func TestBoard_SANToUCI_KnightF3(t *testing.T) {
	b := NewBoard()
	uci, err := b.SANToUCI("Nf3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uci != "g1f3" {
		t.Errorf("got %q, want %q", uci, "g1f3")
	}
}

func TestBoard_SANToUCI_CastlingKingside(t *testing.T) {
	b := NewBoard()
	uci, _ := b.SANToUCI("O-O")
	if uci != "e1g1" {
		t.Errorf("white kingside castling: got %q, want e1g1", uci)
	}
}

func TestBoard_SANToUCI_CastlingQueenside(t *testing.T) {
	b := NewBoard()
	uci, _ := b.SANToUCI("O-O-O")
	if uci != "e1c1" {
		t.Errorf("white queenside castling: got %q, want e1c1", uci)
	}
}

func TestBoard_ApplyUCI_And_BlackMove(t *testing.T) {
	b := NewBoard()
	// White plays e4
	if err := b.ApplyUCI("e2e4"); err != nil {
		t.Fatalf("ApplyUCI e2e4: %v", err)
	}
	if b.turn != ColorBlack {
		t.Error("turn should be Black after White's move")
	}
	// Black plays e5
	uci, err := b.SANToUCI("e5")
	if err != nil {
		t.Fatalf("SANToUCI e5: %v", err)
	}
	if uci != "e7e5" {
		t.Errorf("got %q, want e7e5", uci)
	}
}

func TestBoard_SANToUCI_CheckAnnotationStripped(t *testing.T) {
	b := NewBoard()
	// Knight to f3 with check annotation should still work.
	uci, err := b.SANToUCI("Nf3+")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uci != "g1f3" {
		t.Errorf("got %q, want g1f3", uci)
	}
}

func TestParseScore_CP(t *testing.T) {
	line := "info depth 18 seldepth 22 multipv 1 score cp 35 nodes 1234 nps 500000 time 100 pv e2e4"
	e := parseScore(line)
	if e.Centipawns != 35 || e.IsMate {
		t.Errorf("got %+v, want Centipawns=35 IsMate=false", e)
	}
}

func TestParseScore_Mate(t *testing.T) {
	line := "info depth 5 score mate 2 nodes 100 pv e2e4"
	e := parseScore(line)
	if !e.IsMate || e.MateIn != 2 {
		t.Errorf("got %+v, want IsMate=true MateIn=2", e)
	}
	if e.Centipawns != 9998 {
		t.Errorf("mate 2 centipawns: got %d, want 9998", e.Centipawns)
	}
}
