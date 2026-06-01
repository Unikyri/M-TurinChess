package parser

import (
	"fmt"
)

type Parser struct {
	lexer *Lexer
	curr  Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curr = p.lexer.NextToken()
}

func (p *Parser) parseError(msg string) error {
	return fmt.Errorf("syntax error at line %d, col %d: %s (got %s %q)", p.curr.Line, p.curr.Col, msg, p.curr.Type, p.curr.Value)
}

func (p *Parser) Parse() error {
	// Parse Board
	if err := p.parseBoard(); err != nil {
		return err
	}

	// Parse TK_PIPE
	if p.curr.Type != TK_PIPE {
		return p.parseError("expected '|'")
	}
	p.nextToken()

	// Parse Turn
	if err := p.parseTurn(); err != nil {
		return err
	}

	// Parse TK_PIPE
	if p.curr.Type != TK_PIPE {
		return p.parseError("expected '|'")
	}
	p.nextToken()

	// Parse Castling
	if err := p.parseCastling(); err != nil {
		return err
	}

	// Parse TK_PIPE
	if p.curr.Type != TK_PIPE {
		return p.parseError("expected '|'")
	}
	p.nextToken()

	// Parse EnPassant
	if err := p.parseEnPassant(); err != nil {
		return err
	}

	// Parse TK_EOF
	if p.curr.Type != TK_EOF {
		return p.parseError("expected end of file")
	}

	return nil
}

func (p *Parser) parseBoard() error {
	rowCount := 0
	for {
		if p.curr.Type == TK_ERROR {
			return p.parseError(p.curr.Value)
		}
		// Parse a Row
		sum := 0
		hasSymbols := false
		for p.curr.Type == TK_PIEZA || p.curr.Type == TK_NUMERO {
			hasSymbols = true
			if p.curr.Type == TK_PIEZA {
				sum += 1
			} else if p.curr.Type == TK_NUMERO {
				val := 0
				_, err := fmt.Sscanf(p.curr.Value, "%d", &val)
				if err != nil {
					return p.parseError("invalid number value")
				}
				sum += val
			}
			p.nextToken()
		}

		if !hasSymbols {
			return p.parseError("empty row in board")
		}

		if sum != 8 {
			return fmt.Errorf("row %d sum is %d, expected exactly 8", rowCount+1, sum)
		}

		rowCount++

		if p.curr.Type == TK_SLASH {
			p.nextToken()
		} else {
			break
		}
	}

	if rowCount != 8 {
		return fmt.Errorf("expected exactly 8 rows in board, got %d", rowCount)
	}

	return nil
}

func (p *Parser) parseTurn() error {
	if p.curr.Type == TK_ERROR {
		return p.parseError(p.curr.Value)
	}
	if p.curr.Type != TK_TURNO {
		return p.parseError("expected active color 'w' or 'b'")
	}
	p.nextToken()
	return nil
}

func (p *Parser) parseCastling() error {
	if p.curr.Type == TK_ERROR {
		return p.parseError(p.curr.Value)
	}
	if p.curr.Type != TK_ENROQUE {
		return p.parseError("expected castling rights (subset of KQkq or -)")
	}

	val := p.curr.Value
	if val == "-" {
		p.nextToken()
		return nil
	}

	// We can have a sequence of K, Q, k, q tokens
	seen := make(map[rune]bool)
	for p.curr.Type == TK_ENROQUE {
		charVal := rune(p.curr.Value[0])
		if charVal == '-' {
			return p.parseError("invalid castling option '-' mixed with rights")
		}
		if seen[charVal] {
			return fmt.Errorf("duplicate castling right: %c", charVal)
		}
		seen[charVal] = true
		p.nextToken()
	}

	return nil
}

func (p *Parser) parseEnPassant() error {
	if p.curr.Type == TK_ERROR {
		return p.parseError(p.curr.Value)
	}
	if p.curr.Type != TK_ENPASSANT {
		return p.parseError("expected en passant square (e.g. e3 or -)")
	}
	p.nextToken()
	return nil
}
