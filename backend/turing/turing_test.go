package turing

import (
	"testing"
)

func TestTuringMachine_Acceptance(t *testing.T) {
	// A repeating position sequence
	// Position A repeats 3 times.
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-", // A (1st)
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR|b|KQkq|e3",  // B
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-", // A (2nd)
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR|b|KQkq|e3",  // B
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-", // A (3rd)
	}

	history := fens[:4]
	current := fens[4]

	result := Simulate(history, current)
	if !result.Accepted {
		t.Errorf("Expected simulation to accept triple repetition, but it was rejected (final state: %s)", result.FinalState)
	}

	// Verify that there are three '1's in Tape 3
	tape3Content := make(map[int]string)
	for _, step := range result.Steps {
		tape3Content[step.Tape3.Head] = step.Write3
	}
	onesCount := 0
	for _, cell := range tape3Content {
		if cell == "1" {
			onesCount++
		}
	}
	if onesCount < 3 {
		t.Errorf("Expected at least three '1's on Tape 3, got %d", onesCount)
	}
}

func TestTuringMachine_Rejection(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR|w|KQkq|-", // A
		"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR|b|KQkq|e3",  // B
		"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R|w|KQkq|-",  // C
	}

	history := fens[:2]
	current := fens[2]

	result := Simulate(history, current)
	if result.Accepted {
		t.Errorf("Expected simulation to reject, but it accepted")
	}
}
