package chess

// Symbol represents a classified move in the Turing Machine alphabet.
// M = Módulo (0-10 CP Loss), E = Estándar (11-50), H = Humano (>50).
type Symbol = string

const (
	SymbolM Symbol = "M"
	SymbolE Symbol = "E"
	SymbolH Symbol = "H"
)

// LexerResult holds the classified tape and the subset of moves that were analyzed.
type LexerResult struct {
	Tape          []Symbol
	AnalyzedMoves []MoveAnalysis
}

// openingDepth returns the number of player moves considered "opening theory"
// for a given Elo. Moves within this range that are perfectly played (M) are
// classified as E instead — they are expected to be correct from book knowledge.
func openingDepth(elo int) int {
	switch {
	case elo >= 2200:
		return 6
	case elo >= 1500:
		return 10
	default: // 0–1499 and unknown (0)
		return 15
	}
}

// classifyCP converts a centipawn loss value into the formal alphabet symbol.
func classifyCP(cpLoss int) Symbol {
	switch {
	case cpLoss <= 10:
		return SymbolM
	case cpLoss <= 50:
		return SymbolE
	default:
		return SymbolH
	}
}

// ApplyLexer classifies ALL player moves. Moves inside the opening window that
// are classified as M (0–10 CP Loss) are downgraded to E (Standard) because
// book/theory moves are expected to be optimal and should not raise suspicion.
// This means the full game tape is preserved — no moves are silently dropped.
func ApplyLexer(moves []MoveAnalysis, playerElo int) LexerResult {
	depth := openingDepth(playerElo)

	tape := make([]Symbol, 0, len(moves))
	for i, m := range moves {
		sym := classifyCP(m.CPLoss)
		// Opening phase: a "perfect" move is natural (book theory), not suspicious.
		if i < depth && sym == SymbolM {
			sym = SymbolE
		}
		tape = append(tape, sym)
	}

	return LexerResult{
		Tape:          tape,
		AnalyzedMoves: moves, // ALL moves included
	}
}

