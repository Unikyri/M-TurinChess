package translator

import (
	"fmt"
	"strings"
	"github.com/corentings/chess/v2"
)

// TranslatePGN takes a PGN string and returns a slice of simplified FEN strings,
// one for each position in the game (including the starting position at index 0).
func TranslatePGN(pgnStr string) ([]string, error) {
	// Replace short and long castling 0s with Os.
	// We replace 0-0-0 first so that it doesn't get partially replaced by 0-0.
	pgnStr = strings.ReplaceAll(pgnStr, "0-0-0", "O-O-O")
	pgnStr = strings.ReplaceAll(pgnStr, "0-0", "O-O")

	// If the PGN string is empty, we cannot parse it.
	if strings.TrimSpace(pgnStr) == "" {
		return nil, fmt.Errorf("empty PGN string")
	}

	pgnReader := strings.NewReader(pgnStr)
	pgnOpt, err := chess.PGN(pgnReader)
	if err != nil {
		return nil, fmt.Errorf("error parsing PGN: %w", err)
	}

	game := chess.NewGame(pgnOpt)
	positions := game.Positions()

	simplifiedFENs := make([]string, len(positions))
	for i, pos := range positions {
		fen := pos.String()
		simplified, err := SimplifyFEN(fen)
		if err != nil {
			return nil, err
		}
		simplifiedFENs[i] = simplified
	}

	return simplifiedFENs, nil
}

// SimplifyFEN converts a standard FEN string into the simplified version:
// [Board Representation]|[Active Color]|[Castling Rights]|[En Passant Square]
func SimplifyFEN(fen string) (string, error) {
	parts := strings.Split(fen, " ")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid FEN string: %s", fen)
	}
	return parts[0] + "|" + parts[1] + "|" + parts[2] + "|" + parts[3], nil
}
