package chess

import (
	"fmt"
	"strings"
)

// Color represents the side to move.
type Color int

const (
	ColorWhite Color = iota
	ColorBlack
)

// PieceType identifies a chess piece.
type PieceType int

const (
	Empty PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// Piece holds a type and color.
type Piece struct {
	Type  PieceType
	Color Color
}

var noPiece = Piece{}

// Board is a minimal chess board used only for SAN→UCI conversion.
// It does NOT validate legality (checks, pins) — Stockfish handles that.
type Board struct {
	squares   [64]Piece
	turn      Color
	enPassant int    // target square index; -1 if none
	castling  [4]bool // [WK, WQ, BK, BQ]
}

// NewBoard returns a Board in the standard starting position.
func NewBoard() *Board {
	b := &Board{enPassant: -1, castling: [4]bool{true, true, true, true}}
	b.squares[0] = Piece{Rook, ColorWhite}
	b.squares[1] = Piece{Knight, ColorWhite}
	b.squares[2] = Piece{Bishop, ColorWhite}
	b.squares[3] = Piece{Queen, ColorWhite}
	b.squares[4] = Piece{King, ColorWhite}
	b.squares[5] = Piece{Bishop, ColorWhite}
	b.squares[6] = Piece{Knight, ColorWhite}
	b.squares[7] = Piece{Rook, ColorWhite}
	for i := 8; i < 16; i++ {
		b.squares[i] = Piece{Pawn, ColorWhite}
	}
	for i := 48; i < 56; i++ {
		b.squares[i] = Piece{Pawn, ColorBlack}
	}
	b.squares[56] = Piece{Rook, ColorBlack}
	b.squares[57] = Piece{Knight, ColorBlack}
	b.squares[58] = Piece{Bishop, ColorBlack}
	b.squares[59] = Piece{Queen, ColorBlack}
	b.squares[60] = Piece{King, ColorBlack}
	b.squares[61] = Piece{Bishop, ColorBlack}
	b.squares[62] = Piece{Knight, ColorBlack}
	b.squares[63] = Piece{Rook, ColorBlack}
	return b
}

// --- Square helpers ---

func sqFile(sq int) int { return sq % 8 }
func sqRank(sq int) int { return sq / 8 }
func sqMake(file, rank int) int { return rank*8 + file }

func sqParse(s string) (int, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("invalid square %q", s)
	}
	f, r := int(s[0]-'a'), int(s[1]-'1')
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return 0, fmt.Errorf("square out of range %q", s)
	}
	return sqMake(f, r), nil
}

func sqStr(sq int) string {
	return fmt.Sprintf("%c%c", 'a'+sqFile(sq), '1'+sqRank(sq))
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// ApplyUCI applies a UCI move ("e2e4", "g1f3", "e7e8q") and advances the turn.
func (b *Board) ApplyUCI(move string) error {
	if len(move) < 4 {
		return fmt.Errorf("invalid UCI move %q", move)
	}
	from, err := sqParse(move[0:2])
	if err != nil {
		return err
	}
	to, err := sqParse(move[2:4])
	if err != nil {
		return err
	}
	p := b.squares[from]
	if p.Type == Empty {
		return fmt.Errorf("no piece at %s", move[0:2])
	}
	prevEP := b.enPassant
	b.enPassant = -1

	switch p.Type {
	case King:
		// Castling: king moves 2 files.
		if iabs(sqFile(to)-sqFile(from)) == 2 {
			rank := sqRank(from)
			if sqFile(to) == 6 { // kingside
				b.squares[sqMake(5, rank)] = b.squares[sqMake(7, rank)]
				b.squares[sqMake(7, rank)] = noPiece
			} else { // queenside
				b.squares[sqMake(3, rank)] = b.squares[sqMake(0, rank)]
				b.squares[sqMake(0, rank)] = noPiece
			}
		}
		if p.Color == ColorWhite {
			b.castling[0], b.castling[1] = false, false
		} else {
			b.castling[2], b.castling[3] = false, false
		}
	case Rook:
		if p.Color == ColorWhite {
			if from == 7 {
				b.castling[0] = false
			} else if from == 0 {
				b.castling[1] = false
			}
		} else {
			if from == 63 {
				b.castling[2] = false
			} else if from == 56 {
				b.castling[3] = false
			}
		}
	case Pawn:
		// En-passant capture.
		if to == prevEP {
			b.squares[sqMake(sqFile(to), sqRank(from))] = noPiece
		}
		// Double push: set EP square.
		if iabs(sqRank(to)-sqRank(from)) == 2 {
			b.enPassant = sqMake(sqFile(from), (sqRank(from)+sqRank(to))/2)
		}
		// Promotion.
		if len(move) == 5 {
			pt := Queen
			switch move[4] {
			case 'n':
				pt = Knight
			case 'b':
				pt = Bishop
			case 'r':
				pt = Rook
			}
			b.squares[to] = Piece{pt, p.Color}
			b.squares[from] = noPiece
			b.turn = 1 - b.turn
			return nil
		}
	}
	b.squares[to] = b.squares[from]
	b.squares[from] = noPiece
	b.turn = 1 - b.turn
	return nil
}

// SANToUCI converts a SAN move to UCI notation using the current board state.
func (b *Board) SANToUCI(san string) (string, error) {
	san = strings.TrimRight(san, "+#!?")

	// Castling.
	if san == "O-O-O" {
		if b.turn == ColorWhite {
			return "e1c1", nil
		}
		return "e8c8", nil
	}
	if san == "O-O" {
		if b.turn == ColorWhite {
			return "e1g1", nil
		}
		return "e8g8", nil
	}

	// Promotion suffix.
	promo := ""
	if idx := strings.Index(san, "="); idx >= 0 {
		promo = strings.ToLower(string(san[idx+1]))
		san = san[:idx]
	}

	// Piece type.
	pt := Pawn
	rest := san
	if len(rest) > 0 && rest[0] >= 'A' && rest[0] <= 'Z' {
		switch rest[0] {
		case 'N':
			pt = Knight
		case 'B':
			pt = Bishop
		case 'R':
			pt = Rook
		case 'Q':
			pt = Queen
		case 'K':
			pt = King
		}
		rest = rest[1:]
	}

	// Remove capture marker.
	rest = strings.ReplaceAll(rest, "x", "")

	if len(rest) < 2 {
		return "", fmt.Errorf("cannot parse SAN %q", san)
	}
	destStr := rest[len(rest)-2:]
	disambig := rest[:len(rest)-2]

	toSq, err := sqParse(destStr)
	if err != nil {
		return "", fmt.Errorf("bad destination in %q: %v", san, err)
	}

	fromSq, err := b.findSource(pt, toSq, disambig)
	if err != nil {
		return "", err
	}
	return sqStr(fromSq) + sqStr(toSq) + promo, nil
}

func (b *Board) findSource(pt PieceType, toSq int, disambig string) (int, error) {
	var cands []int
	for sq := 0; sq < 64; sq++ {
		p := b.squares[sq]
		if p.Type != pt || p.Color != b.turn {
			continue
		}
		if b.canReach(sq, toSq, pt) {
			cands = append(cands, sq)
		}
	}
	if len(cands) == 0 {
		return 0, fmt.Errorf("no piece can reach %s", sqStr(toSq))
	}
	if len(cands) == 1 {
		return cands[0], nil
	}
	if disambig == "" {
		return 0, fmt.Errorf("ambiguous move to %s", sqStr(toSq))
	}
	for _, sq := range cands {
		f := string(rune('a' + sqFile(sq)))
		r := string(rune('1' + sqRank(sq)))
		if disambig == f || disambig == r || disambig == sqStr(sq) {
			return sq, nil
		}
	}
	return 0, fmt.Errorf("disambiguation %q did not resolve move to %s", disambig, sqStr(toSq))
}

func (b *Board) canReach(from, to int, pt PieceType) bool {
	// Cannot capture own piece.
	if b.squares[to].Type != Empty && b.squares[to].Color == b.turn {
		return false
	}
	df := sqFile(to) - sqFile(from)
	dr := sqRank(to) - sqRank(from)

	switch pt {
	case Knight:
		return (iabs(df) == 2 && iabs(dr) == 1) || (iabs(df) == 1 && iabs(dr) == 2)
	case King:
		return iabs(df) <= 1 && iabs(dr) <= 1 && (df != 0 || dr != 0)
	case Bishop:
		return iabs(df) == iabs(dr) && df != 0 && !b.blocked(from, to, isign(df), isign(dr))
	case Rook:
		return (df == 0) != (dr == 0) && !b.blocked(from, to, isign(df), isign(dr))
	case Queen:
		if df == 0 && dr == 0 {
			return false
		}
		if df != 0 && dr != 0 && iabs(df) != iabs(dr) {
			return false
		}
		return !b.blocked(from, to, isign(df), isign(dr))
	case Pawn:
		return b.pawnReach(from, to, df, dr)
	}
	return false
}

// blocked returns true if any square between from and to (exclusive) is occupied.
func (b *Board) blocked(from, to, df, dr int) bool {
	f, r := sqFile(from)+df, sqRank(from)+dr
	for {
		sq := sqMake(f, r)
		if sq == to {
			break
		}
		if b.squares[sq].Type != Empty {
			return true
		}
		f += df
		r += dr
	}
	return false
}

func (b *Board) pawnReach(from, to, df, dr int) bool {
	fwd, startRank := 1, 1
	if b.turn == ColorBlack {
		fwd, startRank = -1, 6
	}
	fromRank := sqRank(from)
	// Forward push.
	if df == 0 && dr == fwd && b.squares[to].Type == Empty {
		return true
	}
	// Double push.
	if df == 0 && dr == 2*fwd && fromRank == startRank {
		mid := sqMake(sqFile(from), fromRank+fwd)
		return b.squares[mid].Type == Empty && b.squares[to].Type == Empty
	}
	// Capture or en passant.
	if iabs(df) == 1 && dr == fwd {
		if b.squares[to].Type != Empty && b.squares[to].Color != b.turn {
			return true
		}
		if to == b.enPassant {
			return true
		}
	}
	return false
}
