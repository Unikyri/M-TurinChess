package chess

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// EvalResult holds a Stockfish evaluation.
type EvalResult struct {
	Centipawns int
	IsMate     bool
	MateIn     int // positive = player wins, negative = player loses
}

// MoveAnalysis contains the analysis of a single move.
type MoveAnalysis struct {
	MoveNumber int
	SAN        string
	UCIMove    string
	CPLoss     int
	BestMove   string // UCI notation of the optimal move
	EvalBefore EvalResult
	EvalAfter  EvalResult
}

// StockfishClient manages a long-running Stockfish subprocess.
type StockfishClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
}

// NewStockfishClient starts a Stockfish subprocess and performs the UCI handshake.
func NewStockfishClient(binaryPath string) (*StockfishClient, error) {
	cmd := exec.Command(binaryPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stockfish stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stockfish stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("stockfish binary not found at %q: %w", binaryPath, err)
	}

	sf := &StockfishClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}

	if err := sf.handshake(); err != nil {
		_ = cmd.Process.Kill()
		return nil, err
	}
	return sf, nil
}

func (sf *StockfishClient) send(cmd string) error {
	_, err := fmt.Fprintln(sf.stdin, cmd)
	return err
}

func (sf *StockfishClient) readLine() (string, error) {
	line, err := sf.stdout.ReadString('\n')
	return strings.TrimSpace(line), err
}

func (sf *StockfishClient) handshake() error {
	if err := sf.send("uci"); err != nil {
		return err
	}
	// Wait for uciok with a simple timeout loop.
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		line, err := sf.readLine()
		if err != nil {
			return fmt.Errorf("reading uciok: %w", err)
		}
		if line == "uciok" {
			break
		}
	}
	if err := sf.send("isready"); err != nil {
		return err
	}
	for {
		line, err := sf.readLine()
		if err != nil {
			return fmt.Errorf("reading readyok: %w", err)
		}
		if line == "readyok" {
			return nil
		}
	}
}

// Close sends quit to Stockfish and waits for the process to exit.
func (sf *StockfishClient) Close() {
	_ = sf.send("quit")
	_ = sf.cmd.Wait()
}

// Evaluate evaluates the position reached after playing the given UCI move list.
// Returns the score (from the side TO MOVE's perspective) and the best move.
func (sf *StockfishClient) Evaluate(uciMoves []string, depth int) (EvalResult, string, error) {
	posCmd := "position startpos"
	if len(uciMoves) > 0 {
		posCmd += " moves " + strings.Join(uciMoves, " ")
	}
	if err := sf.send(posCmd); err != nil {
		return EvalResult{}, "", err
	}
	if err := sf.send(fmt.Sprintf("go depth %d", depth)); err != nil {
		return EvalResult{}, "", err
	}
	return sf.readUntilBestMove()
}

// readUntilBestMove reads Stockfish output until "bestmove", returning the
// last score seen and the best move.
func (sf *StockfishClient) readUntilBestMove() (EvalResult, string, error) {
	var lastEval EvalResult
	var bestMove string

	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		line, err := sf.readLine()
		if err != nil {
			return lastEval, bestMove, fmt.Errorf("reading stockfish output: %w", err)
		}
		if strings.HasPrefix(line, "info") && strings.Contains(line, "score") {
			lastEval = parseScore(line)
		}
		if strings.HasPrefix(line, "bestmove") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				bestMove = fields[1]
			}
			return lastEval, bestMove, nil
		}
	}
	return lastEval, bestMove, fmt.Errorf("timeout waiting for bestmove")
}

// parseScore extracts the evaluation from an "info ... score ..." line.
func parseScore(line string) EvalResult {
	fields := strings.Fields(line)
	for i, f := range fields {
		if f == "score" && i+2 < len(fields) {
			kind := fields[i+1]
			val, err := strconv.Atoi(fields[i+2])
			if err != nil {
				return EvalResult{}
			}
			if kind == "cp" {
				return EvalResult{Centipawns: val}
			}
			if kind == "mate" {
				cp := 10000 - iabs(val)
				if val < 0 {
					cp = -cp
				}
				return EvalResult{IsMate: true, MateIn: val, Centipawns: cp}
			}
		}
	}
	return EvalResult{}
}

// AnalyzeGame evaluates every move of the specified color and returns CP Loss per move.
// uciMoves is the full game in UCI notation (all moves, both colors).
func AnalyzeGame(sf *StockfishClient, sanMoves []string, playerColor Color, depth int) ([]MoveAnalysis, error) {
	board := NewBoard()
	uciMoves := make([]string, 0, len(sanMoves))

	var results []MoveAnalysis
	moveNum := 1

	for i, san := range sanMoves {
		// Determine whose turn it is (i=0 is White's first move).
		currentColor := Color(i % 2) // ColorWhite=0, ColorBlack=1

		uci, err := board.SANToUCI(san)
		if err != nil {
			// Skip moves we can't parse rather than aborting the whole game.
			_ = board.ApplyUCI("0000") // null move fallback
			uciMoves = append(uciMoves, "0000")
			if i%2 == 1 {
				moveNum++
			}
			continue
		}

		if currentColor == playerColor {
			// Evaluate position BEFORE this move (from player's perspective).
			evalBefore, bestMove, err := sf.Evaluate(uciMoves, depth)
			if err != nil {
				return nil, fmt.Errorf("evaluate before move %d: %w", moveNum, err)
			}

			// Apply the move.
			_ = board.ApplyUCI(uci)
			uciMoves = append(uciMoves, uci)

			// Evaluate position AFTER this move (from opponent's perspective).
			evalAfter, _, err := sf.Evaluate(uciMoves, depth)
			if err != nil {
				return nil, fmt.Errorf("evaluate after move %d: %w", moveNum, err)
			}

			// CP Loss = max(0, eval_before + eval_after)
			// eval_before: player's perspective; eval_after: opponent's perspective.
			// -eval_after converts opponent perspective back to player perspective.
			cpLoss := evalBefore.Centipawns - (-evalAfter.Centipawns)
			if cpLoss < 0 {
				cpLoss = 0
			}

			results = append(results, MoveAnalysis{
				MoveNumber: moveNum,
				SAN:        san,
				UCIMove:    uci,
				CPLoss:     cpLoss,
				BestMove:   bestMove,
				EvalBefore: evalBefore,
				EvalAfter:  evalAfter,
			})
		} else {
			// Opponent move: just apply it.
			_ = board.ApplyUCI(uci)
			uciMoves = append(uciMoves, uci)
		}

		if i%2 == 1 {
			moveNum++
		}
	}
	return results, nil
}
