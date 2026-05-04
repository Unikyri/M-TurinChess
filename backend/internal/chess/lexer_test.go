package chess

import "testing"

// buildMoves builds a slice of MoveAnalysis with given CP losses for testing.
func buildMoves(cpLosses []int) []MoveAnalysis {
	moves := make([]MoveAnalysis, len(cpLosses))
	for i, cp := range cpLosses {
		moves[i] = MoveAnalysis{MoveNumber: i + 1, CPLoss: cp}
	}
	return moves
}

// --- Classification tests (spec-lexer.md acceptance criteria) ---

func TestClassifyCP_Boundary10(t *testing.T) {
	if classifyCP(10) != SymbolM {
		t.Errorf("CP 10 should be M")
	}
}

func TestClassifyCP_Boundary11(t *testing.T) {
	if classifyCP(11) != SymbolE {
		t.Errorf("CP 11 should be E")
	}
}

func TestClassifyCP_Boundary50(t *testing.T) {
	if classifyCP(50) != SymbolE {
		t.Errorf("CP 50 should be E")
	}
}

func TestClassifyCP_Boundary51(t *testing.T) {
	if classifyCP(51) != SymbolH {
		t.Errorf("CP 51 should be H")
	}
}

// --- Opening filter tests ---

func TestOpeningSkip_Elo1200(t *testing.T) {
	if openingSkipCount(1200) != 15 {
		t.Errorf("Elo 1200: want skip 15, got %d", openingSkipCount(1200))
	}
}

func TestOpeningSkip_Elo1800(t *testing.T) {
	if openingSkipCount(1800) != 10 {
		t.Errorf("Elo 1800: want skip 10, got %d", openingSkipCount(1800))
	}
}

func TestOpeningSkip_Elo2500(t *testing.T) {
	if openingSkipCount(2500) != 6 {
		t.Errorf("Elo 2500: want skip 6, got %d", openingSkipCount(2500))
	}
}

func TestOpeningSkip_EloUnknown(t *testing.T) {
	if openingSkipCount(0) != 15 {
		t.Errorf("Elo 0 (unknown): want skip 15, got %d", openingSkipCount(0))
	}
}

// --- ApplyLexer integration tests ---

func TestApplyLexer_Basic(t *testing.T) {
	// 20 moves for a 1500-2199 Elo player: skip first 10, analyze remaining 10.
	cpLosses := make([]int, 20)
	for i := range cpLosses {
		cpLosses[i] = 5 // all M
	}
	result := ApplyLexer(buildMoves(cpLosses), 1800)
	if len(result.Tape) != 10 {
		t.Errorf("tape length: got %d, want 10", len(result.Tape))
	}
	for i, s := range result.Tape {
		if s != SymbolM {
			t.Errorf("tape[%d]: got %q, want M", i, s)
		}
	}
}

func TestApplyLexer_MixedSymbols(t *testing.T) {
	// Elo 2500 → skip 6 moves. Test 10 moves: first 6 skipped, last 4 classified.
	cpLosses := []int{0, 0, 0, 0, 0, 0, 5, 30, 60, 10}
	result := ApplyLexer(buildMoves(cpLosses), 2500)

	want := []Symbol{SymbolM, SymbolE, SymbolH, SymbolM}
	if len(result.Tape) != len(want) {
		t.Fatalf("tape length: got %d, want %d (tape: %v)", len(result.Tape), len(want), result.Tape)
	}
	for i, s := range want {
		if result.Tape[i] != s {
			t.Errorf("tape[%d]: got %q, want %q", i, result.Tape[i], s)
		}
	}
}

func TestApplyLexer_FewerMovesThanFilter(t *testing.T) {
	// Only 5 moves but Elo 0 needs to skip 15 → empty tape.
	cpLosses := []int{1, 2, 3, 4, 5}
	result := ApplyLexer(buildMoves(cpLosses), 0)
	if len(result.Tape) != 0 {
		t.Errorf("tape should be empty, got %v", result.Tape)
	}
	if len(result.AnalyzedMoves) != 0 {
		t.Errorf("analyzed moves should be empty")
	}
}

func TestApplyLexer_AnalyzedMovesSlice(t *testing.T) {
	// Verify AnalyzedMoves is the post-filter slice.
	cpLosses := []int{5, 5, 5, 5, 5, 5, 25} // 7 moves, Elo 2500 skips 6
	result := ApplyLexer(buildMoves(cpLosses), 2500)
	if len(result.AnalyzedMoves) != 1 {
		t.Errorf("AnalyzedMoves length: got %d, want 1", len(result.AnalyzedMoves))
	}
	if result.AnalyzedMoves[0].CPLoss != 25 {
		t.Errorf("AnalyzedMoves[0].CPLoss: got %d, want 25", result.AnalyzedMoves[0].CPLoss)
	}
}
