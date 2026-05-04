package chess

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Game represents a parsed chess game from a PGN file.
type Game struct {
	White    string
	Black    string
	WhiteElo int    // 0 if not present in the PGN
	BlackElo int    // 0 if not present in the PGN
	Event    string
	Date     string
	Result   string   // "1-0" | "0-1" | "1/2-1/2" | "*"
	Moves    []string // SAN notation, main line only, chronological order
}

// ParsePGN parses a PGN byte slice and returns the first game found.
// Only the main line moves are returned; variations and comments are ignored.
func ParsePGN(data []byte) (Game, error) {
	if len(data) == 0 {
		return Game{}, errors.New("empty pgn file")
	}

	// Normalize line endings.
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	headerSection, movetextSection := splitSections(content)

	var g Game
	parseHeaders(headerSection, &g)
	g.Moves = parseMovetext(movetextSection)

	return g, nil
}

// splitSections separates the PGN into header tags and movetext.
// Headers are lines starting with '['; everything else is movetext.
func splitSections(content string) (string, string) {
	lines := strings.Split(content, "\n")

	var headerLines, movetextLines []string
	headersDone := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !headersDone {
			if strings.HasPrefix(trimmed, "[") {
				headerLines = append(headerLines, line)
			} else if trimmed != "" {
				// First non-tag, non-blank line signals start of movetext.
				headersDone = true
				movetextLines = append(movetextLines, line)
			}
			// Blank lines between headers are ignored.
		} else {
			movetextLines = append(movetextLines, line)
		}
	}

	return strings.Join(headerLines, "\n"), strings.Join(movetextLines, "\n")
}

// parseHeaders extracts known metadata fields from the PGN tag section.
func parseHeaders(section string, g *Game) {
	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
			continue
		}

		// Format: [Key "Value"]
		inner := line[1 : len(line)-1]
		spaceIdx := strings.Index(inner, " ")
		if spaceIdx < 0 {
			continue
		}

		key := strings.TrimSpace(inner[:spaceIdx])
		value := strings.TrimSpace(inner[spaceIdx+1:])

		// Remove surrounding double quotes.
		if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		switch key {
		case "White":
			g.White = value
		case "Black":
			g.Black = value
		case "WhiteElo":
			if n, err := strconv.Atoi(value); err == nil {
				g.WhiteElo = n
			}
		case "BlackElo":
			if n, err := strconv.Atoi(value); err == nil {
				g.BlackElo = n
			}
		case "Event":
			g.Event = value
		case "Date":
			g.Date = value
		case "Result":
			g.Result = value
		}
	}
}

var moveNumRegex = regexp.MustCompile(`\b\d+\.+`)

// parseMovetext extracts only the main line SAN moves from the movetext section.
// It strips: comments {...}, variations (...), NAGs $N, move numbers, and the result token.
func parseMovetext(section string) []string {
	section = removeComments(section)
	section = removeVariations(section)
	section = removeNAGs(section)
	
	// Replace move numbers like "1." or "15..." with spaces to detach them from moves like "1.e4"
	section = moveNumRegex.ReplaceAllString(section, " ")
	
	return extractSANTokens(section)
}

// removeComments removes all {...} comment blocks (supports nesting).
func removeComments(s string) string {
	var b strings.Builder
	depth := 0
	for _, ch := range s {
		switch ch {
		case '{':
			depth++
		case '}':
			if depth > 0 {
				depth--
			}
		default:
			if depth == 0 {
				b.WriteRune(ch)
			}
		}
	}
	return b.String()
}

// removeVariations removes all (...) variation blocks, including nested ones.
func removeVariations(s string) string {
	var b strings.Builder
	depth := 0
	for _, ch := range s {
		switch ch {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		default:
			if depth == 0 {
				b.WriteRune(ch)
			}
		}
	}
	return b.String()
}

// removeNAGs removes Numeric Annotation Glyphs like $1, $2, $21.
func removeNAGs(s string) string {
	var b strings.Builder
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		if runes[i] == '$' {
			i++ // skip '$'
			for i < len(runes) && runes[i] >= '0' && runes[i] <= '9' {
				i++ // skip digits
			}
		} else {
			b.WriteRune(runes[i])
			i++
		}
	}
	return b.String()
}

// gameTerminators are valid PGN result tokens that mark end of game.
var gameTerminators = map[string]bool{
	"1-0": true, "0-1": true, "1/2-1/2": true, "*": true,
}

// isMoveNumber returns true if the token is a move number like "1.", "15.", "1...".
func isMoveNumber(token string) bool {
	stripped := strings.TrimRight(token, ".")
	if stripped == "" || !strings.ContainsRune(token, '.') {
		return false
	}
	for _, ch := range stripped {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// extractSANTokens tokenizes the cleaned movetext and returns only SAN moves.
// Stops at the first game terminator.
func extractSANTokens(s string) []string {
	var moves []string
	for _, token := range strings.Fields(s) {
		if gameTerminators[token] {
			break
		}
		if isMoveNumber(token) || token == "" {
			continue
		}
		moves = append(moves, token)
	}
	return moves
}
