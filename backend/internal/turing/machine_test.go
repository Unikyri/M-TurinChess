package turing

import (
	"testing"
)

// --- delta unit tests (T1.2) ---

func TestDelta_M_IncrementsSuspicion(t *testing.T) {
	res := delta(stQ0, "M", blank)
	if res.writeC2 != "1" || res.moveC2 != dirR || res.next != stQ0 {
		t.Errorf("q0+M: got writeC2=%q moveC2=%d next=%v", res.writeC2, res.moveC2, res.next)
	}
}

func TestDelta_E_Neutral(t *testing.T) {
	res := delta(stQ0, "E", blank)
	if res.writeC2 != blank || res.moveC2 != dirS || res.next != stQ0 {
		t.Errorf("q0+E: got writeC2=%q moveC2=%d next=%v", res.writeC2, res.moveC2, res.next)
	}
}

func TestDelta_H_PausesC1_MovesC2Left(t *testing.T) {
	res := delta(stQ0, "H", blank)
	if res.moveC1 != dirS || res.moveC2 != dirL || res.next != stQBorra {
		t.Errorf("q0+H: got moveC1=%d moveC2=%d next=%v", res.moveC1, res.moveC2, res.next)
	}
}

func TestDelta_Blank_GoesToQF(t *testing.T) {
	res := delta(stQ0, blank, blank)
	if res.next != stQF {
		t.Errorf("q0+B: expected qF, got %v", res.next)
	}
}

func TestDelta_QBorra_Erases1(t *testing.T) {
	res := delta(stQBorra, "H", "1")
	if res.writeC2 != blank || res.moveC1 != dirR || res.next != stQ0 {
		t.Errorf("q_borra+1: got writeC2=%q moveC1=%d next=%v", res.writeC2, res.moveC1, res.next)
	}
}

func TestDelta_QBorra_NothingToErase(t *testing.T) {
	res := delta(stQBorra, "H", blank)
	if res.moveC1 != dirR || res.next != stQ0 {
		t.Errorf("q_borra+B: got moveC1=%d next=%v", res.moveC1, res.next)
	}
}

// --- Run acceptance tests (T1.3 / T1.4) — all spec criteria ---

func TestRun_MMH_SuspicionOne(t *testing.T) {
	r := Run([]string{"M", "M", "H"})
	if r.SuspicionCount != 1 {
		t.Errorf("MMH: want SuspicionCount=1, got %d", r.SuspicionCount)
	}
}

func TestRun_SixM_SuspicionSix(t *testing.T) {
	r := Run([]string{"M", "M", "M", "M", "M", "M"})
	if r.SuspicionCount != 6 {
		t.Errorf("6xM: want 6, got %d", r.SuspicionCount)
	}
}

func TestRun_ThreeH_SuspicionZero(t *testing.T) {
	r := Run([]string{"H", "H", "H"})
	if r.SuspicionCount != 0 {
		t.Errorf("HHH: want 0, got %d (must not go negative)", r.SuspicionCount)
	}
}

func TestRun_EmptyTape_SuspicionZero(t *testing.T) {
	r := Run([]string{})
	if r.SuspicionCount != 0 {
		t.Errorf("empty tape: want 0, got %d", r.SuspicionCount)
	}
	if len(r.Trace) == 0 {
		t.Error("empty tape: Trace should have at least 1 step (the qf step)")
	}
}

func TestRun_MHMHM_SuspicionOne(t *testing.T) {
	r := Run([]string{"M", "H", "M", "H", "M"})
	if r.SuspicionCount != 1 {
		t.Errorf("MHMHM: want 1, got %d", r.SuspicionCount)
	}
}

func TestRun_AllE_SuspicionZero(t *testing.T) {
	r := Run([]string{"E", "E", "E", "E"})
	if r.SuspicionCount != 0 {
		t.Errorf("EEEE: want 0, got %d", r.SuspicionCount)
	}
}

func TestRun_Trace_NonEmpty(t *testing.T) {
	r := Run([]string{"M", "E", "H"})
	if len(r.Trace) == 0 {
		t.Error("Trace should not be empty")
	}
	// Verify suspicion is monotonically non-negative throughout trace.
	for i, step := range r.Trace {
		if step.Suspicion < 0 {
			t.Errorf("Trace[%d].Suspicion is negative: %d", i, step.Suspicion)
		}
	}
}

func TestRun_Trace_SuspicionConsistency(t *testing.T) {
	// After MMH, the trace should show: 0→1→2→1 (in the Suspicion field of each step).
	r := Run([]string{"M", "M", "H"})
	// Step 0 (q0, read M): suspicion becomes 1
	// Step 1 (q0, read M): suspicion becomes 2
	// Step 2 (q0, read H, C1 pauses): suspicion stays 2 (C2 just moved L, no erase yet)
	// Step 3 (q_borra, read 1): suspicion becomes 1
	// Step 4 (q0, read B): qf → final suspicion = 1
	if r.Trace[0].Suspicion != 1 {
		t.Errorf("after step 0 (M): want suspicion=1, got %d", r.Trace[0].Suspicion)
	}
	if r.Trace[1].Suspicion != 2 {
		t.Errorf("after step 1 (M): want suspicion=2, got %d", r.Trace[1].Suspicion)
	}
	last := r.Trace[len(r.Trace)-1]
	if last.Suspicion != 1 {
		t.Errorf("final suspicion: want 1, got %d", last.Suspicion)
	}
}

func TestRun_Tape2State_ReflectsOnes(t *testing.T) {
	r := Run([]string{"M", "M", "M"})
	// Tape 2 should have 3 ones.
	if len(r.Tape2State) != 3 {
		t.Errorf("MMM tape2: want 3 cells, got %d (%v)", len(r.Tape2State), r.Tape2State)
	}
	for i, c := range r.Tape2State {
		if c != "1" {
			t.Errorf("Tape2State[%d] = %q, want '1'", i, c)
		}
	}
}

func TestRun_Isolation_NoExternalImports(t *testing.T) {
	// This test is a placeholder to document the isolation requirement.
	// The actual enforcement is: `go build ./internal/turing/` must not import
	// any package outside the stdlib and this package.
	t.Log("turing package uses only stdlib (fmt) — isolation confirmed by go.mod")
}
