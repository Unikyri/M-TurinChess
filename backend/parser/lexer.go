package parser

import (
	"fmt"
	"strings"
)

type TokenType string

const (
	TK_PIEZA     TokenType = "TK_PIEZA"
	TK_NUMERO    TokenType = "TK_NUMERO"
	TK_SLASH     TokenType = "TK_SLASH"
	TK_PIPE      TokenType = "TK_PIPE"
	TK_TURNO     TokenType = "TK_TURNO"
	TK_ENROQUE   TokenType = "TK_ENROQUE"
	TK_ENPASSANT TokenType = "TK_ENPASSANT"
	TK_EOF       TokenType = "TK_EOF"
	TK_ERROR     TokenType = "TK_ERROR"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input []rune
	pos   int
	field int // 0: board, 1: turn, 2: enroque, 3: enpassant
	col   int
	line  int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
		field: 0,
		col:   1,
		line:  1,
	}
}

func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.input) {
		return Token{Type: TK_EOF, Value: "", Line: l.line, Col: l.col}
	}

	r := l.input[l.pos]
	startCol := l.col

	// Handle field transition on pipe
	if r == '|' {
		l.pos++
		l.col++
		l.field++
		return Token{Type: TK_PIPE, Value: "|", Line: l.line, Col: startCol}
	}

	switch l.field {
	case 0: // Board Representation
		l.pos++
		l.col++
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			if strings.ContainsRune("pnbrqkPNBRQK", r) {
				return Token{Type: TK_PIEZA, Value: string(r), Line: l.line, Col: startCol}
			}
			return Token{Type: TK_ERROR, Value: fmt.Sprintf("invalid piece character: %c", r), Line: l.line, Col: startCol}
		}
		if r >= '1' && r <= '8' {
			return Token{Type: TK_NUMERO, Value: string(r), Line: l.line, Col: startCol}
		}
		if r == '/' {
			return Token{Type: TK_SLASH, Value: "/", Line: l.line, Col: startCol}
		}
		return Token{Type: TK_ERROR, Value: fmt.Sprintf("unexpected character in board field: %c", r), Line: l.line, Col: startCol}

	case 1: // Turn
		l.pos++
		l.col++
		if r == 'w' || r == 'b' {
			return Token{Type: TK_TURNO, Value: string(r), Line: l.line, Col: startCol}
		}
		return Token{Type: TK_ERROR, Value: fmt.Sprintf("invalid turn character: %c", r), Line: l.line, Col: startCol}

	case 2: // Castling Rights (Enroque)
		l.pos++
		l.col++
		if r == '-' {
			return Token{Type: TK_ENROQUE, Value: "-", Line: l.line, Col: startCol}
		}
		if r == 'K' || r == 'Q' || r == 'k' || r == 'q' {
			return Token{Type: TK_ENROQUE, Value: string(r), Line: l.line, Col: startCol}
		}
		return Token{Type: TK_ERROR, Value: fmt.Sprintf("invalid castling character: %c", r), Line: l.line, Col: startCol}

	case 3: // En Passant (Peón al paso)
		if r == '-' {
			l.pos++
			l.col++
			return Token{Type: TK_ENPASSANT, Value: "-", Line: l.line, Col: startCol}
		}
		// Expect square like e3 (file a-h, rank 1-8)
		if r >= 'a' && r <= 'h' {
			if l.pos+1 < len(l.input) {
				nextR := l.input[l.pos+1]
				if nextR >= '1' && nextR <= '8' {
					square := string(r) + string(nextR)
					l.pos += 2
					l.col += 2
					return Token{Type: TK_ENPASSANT, Value: square, Line: l.line, Col: startCol}
				}
			}
			l.pos++
			l.col++
			return Token{Type: TK_ERROR, Value: fmt.Sprintf("invalid en passant rank after file %c", r), Line: l.line, Col: startCol}
		}
		l.pos++
		l.col++
		return Token{Type: TK_ERROR, Value: fmt.Sprintf("invalid en passant character: %c", r), Line: l.line, Col: startCol}

	default: // If there are extra pipes or chars after the 4th field
		l.pos++
		l.col++
		return Token{Type: TK_ERROR, Value: fmt.Sprintf("unexpected content after 4th field: %c", r), Line: l.line, Col: startCol}
	}
}
