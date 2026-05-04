package turing

import "fmt"

// TraceStep records a single step of the Turing Machine execution.
type TraceStep struct {
	Step      int    `json:"step"`
	State     string `json:"state"`
	ReadC1    string `json:"read_c1"`
	ReadC2    string `json:"read_c2"`
	Action    string `json:"action"`
	Suspicion int    `json:"suspicion"`
}

// MTResult holds the outcome of a Turing Machine run.
type MTResult struct {
	SuspicionCount int        `json:"suspicion_count"`
	Tape2State     []string   `json:"tape2_state"`
	Trace          []TraceStep `json:"trace"`
}

// --- Internal state and direction types ---

type tmState int

const (
	stQ0     tmState = iota
	stQBorra
	stQF
)

func (s tmState) String() string {
	switch s {
	case stQ0:
		return "q0"
	case stQBorra:
		return "q_borra"
	default:
		return "qf"
	}
}

type direction int

const (
	dirL direction = -1
	dirS direction = 0
	dirR direction = 1
)

const blank = "B"

// --- Semi-infinite tape ---

// tmTape is a semi-infinite tape that auto-extends rightward.
// The leftmost cell is index 0; the head cannot go past 0 (left boundary guard).
type tmTape struct {
	cells []string
	head  int
}

func (t *tmTape) read() string {
	if t.head >= len(t.cells) || t.cells[t.head] == "" {
		return blank
	}
	return t.cells[t.head]
}

func (t *tmTape) write(sym string) {
	for t.head >= len(t.cells) {
		t.cells = append(t.cells, blank)
	}
	t.cells[t.head] = sym
}

func (t *tmTape) move(d direction) {
	t.head += int(d)
	if t.head < 0 {
		t.head = 0 // left boundary guard — head never goes negative
	}
}

// --- Transition function δ ---

type stepResult struct {
	writeC1 string
	moveC1  direction
	writeC2 string
	moveC2  direction
	next    tmState
}

// delta applies the transition table δ(state, c1, c2) → stepResult.
//
// Updated transition table (from spec-tm-simulator.md):
//
//	q0,     M, B → M, R,  1,  R,  q0
//	q0,     E, B → E, R,  B,  S,  q0
//	q0,     H, B → H, S,  B,  L*, q_borra   ← C1 pauses; C2 moves left
//	q0,     B, * → B, S,  *,  S,  qf
//	q_borra,*, 1 → *, R,  B,  S,  q0        ← erase 1; C1 resumes
//	q_borra,*, B → *, R,  B,  S,  q0        ← nothing to erase; C1 resumes
func delta(st tmState, c1, c2 string) stepResult {
	switch st {
	case stQ0:
		switch c1 {
		case "M":
			return stepResult{"M", dirR, "1", dirR, stQ0}
		case "E":
			return stepResult{"E", dirR, c2, dirS, stQ0}
		case "H":
			// C1 stays (S); C2 moves left to position of the last 1 (or stays at 0 if empty).
			return stepResult{"H", dirS, c2, dirL, stQBorra}
		default: // blank or any unknown symbol → accept
			return stepResult{c1, dirS, c2, dirS, stQF}
		}
	case stQBorra:
		// C1 was paused at H; it now advances R to consume it.
		if c2 == "1" {
			return stepResult{c1, dirR, blank, dirS, stQ0} // erase the 1
		}
		return stepResult{c1, dirR, c2, dirS, stQ0} // count was 0; nothing to erase
	}
	// stQF is absorbing.
	return stepResult{c1, dirS, c2, dirS, stQF}
}

// --- Helpers ---

func countOnes(cells []string) int {
	n := 0
	for _, c := range cells {
		if c == "1" {
			n++
		}
	}
	return n
}

func actionDesc(st tmState, c1, c2 string) string {
	switch {
	case st == stQ0 && c1 == "M":
		return "M: escribe 1 en C2, cabezal C2 →"
	case st == stQ0 && c1 == "E":
		return "E: neutro, sin cambio en C2"
	case st == stQ0 && c1 == "H":
		return "H: pausa C1, mueve C2 ← hacia el último 1"
	case st == stQ0:
		return "B: fin de cinta → qf"
	case st == stQBorra && c2 == "1":
		return "q_borra: borra 1 (sospecha--), reanuda C1"
	case st == stQBorra:
		return fmt.Sprintf("q_borra: C2=%q, nada que borrar (count=0), reanuda C1", c2)
	}
	return "qf: aceptado"
}

// --- Public API ---

// Run simulates the 2-tape Turing Machine on the given input tape.
// Returns the suspicion count, the final state of Tape 2, and a full execution trace.
func Run(input []string) MTResult {
	// Initialise Tape 1: input symbols + explicit blank sentinel.
	c1cells := make([]string, len(input)+1)
	copy(c1cells, input)
	c1cells[len(input)] = blank
	c1 := &tmTape{cells: c1cells}

	// Initialise Tape 2: empty (all blanks on demand).
	c2 := &tmTape{cells: []string{}}

	st := stQ0
	var trace []TraceStep
	stepNum := 0

	for {
		r1 := c1.read()
		r2 := c2.read()
		res := delta(st, r1, r2)

		// Apply writes first, then moves.
		c1.write(res.writeC1)
		c2.write(res.writeC2)
		c1.move(res.moveC1)
		c2.move(res.moveC2)

		suspicion := countOnes(c2.cells)

		trace = append(trace, TraceStep{
			Step:      stepNum,
			State:     st.String(),
			ReadC1:    r1,
			ReadC2:    r2,
			Action:    actionDesc(st, r1, r2),
			Suspicion: suspicion,
		})

		st = res.next
		stepNum++

		if st == stQF {
			break
		}
	}

	// Build the meaningful Tape 2 state (trim trailing blanks for display).
	tape2 := make([]string, len(c2.cells))
	copy(tape2, c2.cells)
	for len(tape2) > 0 && tape2[len(tape2)-1] == blank {
		tape2 = tape2[:len(tape2)-1]
	}
	if tape2 == nil {
		tape2 = []string{}
	}

	return MTResult{
		SuspicionCount: countOnes(c2.cells),
		Tape2State:     tape2,
		Trace:          trace,
	}
}
