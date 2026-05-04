package analysis

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Unikyri/M-TurinChess/backend/internal/chess"
	"github.com/Unikyri/M-TurinChess/backend/internal/db"
	"github.com/Unikyri/M-TurinChess/backend/internal/llm"
	"github.com/Unikyri/M-TurinChess/backend/internal/turing"
)

// AnalysisRequest holds all parameters for a single game analysis.
type AnalysisRequest struct {
	PGNBytes    []byte
	PlayerColor string // "white" | "black"
	EloWhite    int    // 0 means not provided; will fall back to PGN tags
	EloBlack    int    // 0 means not provided; will fall back to PGN tags
	Threshold   int    // suspicion threshold; default 6
	Depth       int    // Stockfish depth; default 18
}

// MoveDetail is a single move's analysis result, suitable for JSON serialization.
type MoveDetail struct {
	MoveNumber     int     `json:"move_number"`
	SAN            string  `json:"san"`
	UCI            string  `json:"uci"`
	CPLoss         int     `json:"cp_loss"`
	Classification string  `json:"classification"`
	BestMove       string  `json:"best_move"`
	LLMFlag        *string `json:"llm_flag"`
}

// AnalysisResult is the full output of an analysis, returned to the API caller.
type AnalysisResult struct {
	ID                 string             `json:"id"`
	AnalyzedAt         time.Time          `json:"analyzed_at"`
	Verdict            string             `json:"verdict"` // "MODULE_DETECTED" | "HUMAN_PLAYER"
	SuspicionCount     int                `json:"suspicion_count"`
	Threshold          int                `json:"threshold"`
	TotalMovesAnalyzed int                `json:"total_moves_analyzed"`
	PlayerColor        string             `json:"player_color"`
	Elo                int                `json:"elo"`
	TapeInput          []string           `json:"tape_input"`
	TapeOutput         []string           `json:"tape_output"`
	MoveDetails        []MoveDetail       `json:"move_details"`
	MTTrace            []turing.TraceStep `json:"mt_trace"`
}

// Analyzer orchestrates the full analysis pipeline.
type Analyzer struct {
	stockfishPath string
	storage       db.Storage
	llmClient     llm.Client
}

// NewAnalyzer creates an Analyzer with the given Stockfish binary path, storage, and optional LLM client.
func NewAnalyzer(stockfishPath string, storage db.Storage, llmClient llm.Client) *Analyzer {
	return &Analyzer{stockfishPath: stockfishPath, storage: storage, llmClient: llmClient}
}

// Analyze runs the complete pipeline: PGN → Stockfish → Lexer → MT → Verdict → DB.
func (a *Analyzer) Analyze(req AnalysisRequest) (AnalysisResult, error) {
	// Apply defaults.
	if req.Threshold <= 0 {
		req.Threshold = 6
	}
	if req.Depth <= 0 {
		req.Depth = 6 // TEMP: lowered for LLM verification (restore to 18)
	}

	// Step 1: Parse PGN.
	game, err := chess.ParsePGN(req.PGNBytes)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("parse pgn: %w", err)
	}

	// Step 2: Determine player color and Elo.
	playerColor := chess.ColorWhite
	playerElo := req.EloWhite
	if req.PlayerColor == "black" {
		playerColor = chess.ColorBlack
		playerElo = req.EloBlack
	}
	// Fall back to PGN-embedded Elo if not supplied in request.
	if playerElo == 0 {
		if req.PlayerColor == "white" {
			playerElo = game.WhiteElo
		} else {
			playerElo = game.BlackElo
		}
	}

	// Step 3: Start Stockfish and evaluate every player move.
	sf, err := chess.NewStockfishClient(a.stockfishPath)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("start stockfish: %w", err)
	}
	defer sf.Close()

	moveAnalyses, err := chess.AnalyzeGame(sf, game.Moves, playerColor, req.Depth)
	if err != nil {
		return AnalysisResult{}, fmt.Errorf("analyze game: %w", err)
	}

	// Step 4: Apply lexer (opening filter + classification).
	lexResult := chess.ApplyLexer(moveAnalyses, playerElo)

	// Step 5: Gather LLM context (Entire game for better endgame understanding)
	var contextParts []string
	for _, am := range lexResult.AnalyzedMoves {
		contextParts = append(contextParts, fmt.Sprintf("%d.%s", am.MoveNumber, am.SAN))
	}
	moveContext := strings.Join(contextParts, " ")

	// Collect all "M" moves to send to LLM in one batch.
	var suspiciousSANs []string
	for i, m := range lexResult.AnalyzedMoves {
		if lexResult.Tape[i] == "M" {
			suspiciousSANs = append(suspiciousSANs, m.SAN)
		}
	}

	// Call LLM for the batch
	llmResults := make(map[string]string)
	if a.llmClient != nil && len(suspiciousSANs) > 0 {
		llmResults = a.llmClient.AnalyzeBatch(suspiciousSANs, playerElo, moveContext)
	}

	// Override Lexer Tape with LLM's superior judgment BEFORE running Turing Machine
	details := make([]MoveDetail, len(lexResult.AnalyzedMoves))
	for i, m := range lexResult.AnalyzedMoves {
		var llmFlag *string
		if lexResult.Tape[i] == "M" {
			if flag, ok := llmResults[m.SAN]; ok {
				flagCopy := flag
				llmFlag = &flagCopy

				// If LLM says it's Human or Standard, downgrade the Tape!
				if strings.HasPrefix(flagCopy, "H") {
					lexResult.Tape[i] = "H"
				} else if strings.HasPrefix(flagCopy, "E") {
					lexResult.Tape[i] = "E"
				}
			}
		}

		details[i] = MoveDetail{
			MoveNumber:     m.MoveNumber,
			SAN:            m.SAN,
			UCI:            m.UCIMove,
			CPLoss:         m.CPLoss,
			Classification: lexResult.Tape[i], // Shows the updated classification!
			BestMove:       m.BestMove,
			LLMFlag:        llmFlag,
		}
	}

	// Step 6: Run Turing Machine on the UPDATED tape!
	tapeCopy := make([]string, len(lexResult.Tape))
	copy(tapeCopy, lexResult.Tape)
	mtResult := turing.Run(tapeCopy)

	// Step 7: Determine verdict based on the LLM-filtered tape.
	verdict := "HUMAN_PLAYER"
	if mtResult.SuspicionCount >= req.Threshold {
		verdict = "MODULE_DETECTED"
	}

	// Step 8: Normalise tape output.
	tapeOut := mtResult.Tape2State
	if tapeOut == nil {
		tapeOut = []string{}
	}

	id := newID()
	result := AnalysisResult{
		ID:                 id,
		AnalyzedAt:         time.Now().UTC(),
		Verdict:            verdict,
		SuspicionCount:     mtResult.SuspicionCount,
		Threshold:          req.Threshold,
		TotalMovesAnalyzed: len(lexResult.AnalyzedMoves),
		PlayerColor:        req.PlayerColor,
		Elo:                playerElo,
		TapeInput:          lexResult.Tape,
		TapeOutput:         tapeOut,
		MoveDetails:        details,
		MTTrace:            mtResult.Trace,
	}

	// Step 9: Persist summary to JSON DB (non-fatal on error).
	record := db.HistoryRecord{
		ID:                 id,
		AnalyzedAt:         result.AnalyzedAt,
		Verdict:            result.Verdict,
		SuspicionCount:     result.SuspicionCount,
		Threshold:          result.Threshold,
		TotalMovesAnalyzed: result.TotalMovesAnalyzed,
		PlayerColor:        result.PlayerColor,
		Elo:                result.Elo,
	}
	if err := a.storage.Append(record); err != nil {
		fmt.Printf("[warn] persist analysis: %v\n", err)
	}

	return result, nil
}

// newID generates a random 16-character hex string for use as a record ID.
func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
