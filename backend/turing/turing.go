package turing

import (
	"strings"
)

type Tape struct {
	Content map[int]string
	Head    int
	MinIdx  int
	MaxIdx  int
}

func NewTape(initialContent string) *Tape {
	t := &Tape{
		Content: make(map[int]string),
		Head:    0,
		MinIdx:  0,
		MaxIdx:  -1,
	}
	for i, char := range initialContent {
		t.Content[i] = string(char)
		t.MaxIdx = i
	}
	if len(initialContent) == 0 {
		t.MinIdx = 0
		t.MaxIdx = 0
		t.Content[0] = "_"
	}
	return t
}

func (t *Tape) Read() string {
	val, ok := t.Content[t.Head]
	if !ok {
		return "_"
	}
	return val
}

func (t *Tape) Write(val string) {
	t.Content[t.Head] = val
	if t.Head < t.MinIdx {
		t.MinIdx = t.Head
	}
	if t.Head > t.MaxIdx {
		t.MaxIdx = t.Head
	}
}

func (t *Tape) Move(dir string) {
	switch dir {
	case "R":
		t.Head++
	case "L":
		t.Head--
	}
}

func (t *Tape) GetCellsAndHead() ([]string, int, int) {
	min := t.MinIdx
	max := t.MaxIdx
	if t.Head < min {
		min = t.Head
	}
	if t.Head > max {
		max = t.Head
	}
	cells := make([]string, max-min+1)
	for i := min; i <= max; i++ {
		val, ok := t.Content[i]
		if !ok {
			cells[i-min] = "_"
		} else {
			cells[i-min] = val
		}
	}
	return cells, t.Head - min, min
}

type TapeState struct {
	Head   int      `json:"head"`
	MinIdx int      `json:"minIdx"`
}

type StepTrace struct {
	Step      int       `json:"step"`
	State     string    `json:"state"`
	Tape1     TapeState `json:"tape1"`
	Tape2     TapeState `json:"tape2"`
	Tape3     TapeState `json:"tape3"`
	Read1     string    `json:"read1"`
	Read2     string    `json:"read2"`
	Read3     string    `json:"read3"`
	Write1    string    `json:"write1"`
	Write2    string    `json:"write2"`
	Write3    string    `json:"write3"`
	Dir1      string    `json:"dir1"`
	Dir2      string    `json:"dir2"`
	Dir3      string    `json:"dir3"`
	NextState string    `json:"nextState"`
}

type SimulationResult struct {
	Steps        []StepTrace `json:"steps,omitempty"`
	FinalState   string      `json:"finalState"`
	Accepted     bool        `json:"accepted"`
	Tape1Initial string      `json:"tape1Initial"`
	Tape2Initial string      `json:"tape2Initial"`
	Tape3Initial string      `json:"tape3Initial"`
}

func Simulate(historyFENs []string, currentFEN string) *SimulationResult {
	var tape1Str string
	if len(historyFENs) > 0 {
		tape1Str = "$" + strings.Join(historyFENs, "$") + "$"
	} else {
		tape1Str = "$"
	}
	tape2Str := "$" + currentFEN + "$"

	tape1 := NewTape(tape1Str)
	tape2 := NewTape(tape2Str)
	tape3 := NewTape("")

	state := "q_init"
	steps := []StepTrace{}
	stepNum := 0
	maxSteps := 50000

	for stepNum < maxSteps {
		r1 := tape1.Read()
		r2 := tape2.Read()
		r3 := tape3.Read()

		t1Min := tape1.MinIdx
		if tape1.Head < t1Min {
			t1Min = tape1.Head
		}
		t1State := TapeState{Head: tape1.Head, MinIdx: t1Min}

		t2Min := tape2.MinIdx
		if tape2.Head < t2Min {
			t2Min = tape2.Head
		}
		t2State := TapeState{Head: tape2.Head, MinIdx: t2Min}

		t3Min := tape3.MinIdx
		if tape3.Head < t3Min {
			t3Min = tape3.Head
		}
		t3State := TapeState{Head: tape3.Head, MinIdx: t3Min}

		w1, w2, w3 := r1, r2, r3
		d1, d2, d3 := "S", "S", "S"
		var nextState string

		switch state {
		case "q_init":
			w3 = "$"
			d3 = "R"
			nextState = "q_init_write1"

		case "q_init_write1":
			w3 = "1"
			d3 = "R"
			nextState = "q_cmp"

		case "q_cmp":
			if r1 == "_" {
				d3 = "L"
				nextState = "q_rewindC3"
			} else if r2 == "_" {
				w3 = "1"
				d3 = "R"
				d2 = "L"
				nextState = "q_rewindC2_skip"
			} else if r1 == r2 {
				d1 = "R"
				d2 = "R"
				nextState = "q_cmp"
			} else {
				d2 = "L"
				nextState = "q_rebobinarC2"
			}

		case "q_rewindC2_skip":
			d2 = "L"
			nextState = "q_rebobinarC2"

		case "q_rebobinarC2":
			if r2 != "$" {
				d2 = "L"
				nextState = "q_rebobinarC2"
			} else {
				nextState = "q_saltarC1"
			}

		case "q_saltarC1":
			if r1 == "_" {
				d3 = "L"
				nextState = "q_rewindC3"
			} else if r1 != "$" {
				d1 = "R"
				nextState = "q_saltarC1"
			} else {
				nextState = "q_cmp"
			}

		case "q_rewindC3":
			if r3 == "1" {
				d3 = "L"
				nextState = "q_rewindC3"
			} else if r3 == "$" {
				d3 = "R"
				nextState = "q_countC3_1"
			} else {
				nextState = "q_reject"
			}

		case "q_countC3_1":
			if r3 == "1" {
				d3 = "R"
				nextState = "q_countC3_2"
			} else {
				nextState = "q_reject"
			}

		case "q_countC3_2":
			if r3 == "1" {
				d3 = "R"
				nextState = "q_countC3_3"
			} else {
				nextState = "q_reject"
			}

		case "q_countC3_3":
			if r3 == "1" {
				nextState = "q_accept"
			} else {
				nextState = "q_reject"
			}

		case "q_accept", "q_reject":
			nextState = state
		}

		steps = append(steps, StepTrace{
			Step:      stepNum,
			State:     state,
			Tape1:     t1State,
			Tape2:     t2State,
			Tape3:     t3State,
			Read1:     r1,
			Read2:     r2,
			Read3:     r3,
			Write1:    w1,
			Write2:    w2,
			Write3:    w3,
			Dir1:      d1,
			Dir2:      d2,
			Dir3:      d3,
			NextState: nextState,
		})

		if state == "q_accept" || state == "q_reject" {
			break
		}

		tape1.Write(w1)
		tape2.Write(w2)
		tape3.Write(w3)

		tape1.Move(d1)
		tape2.Move(d2)
		tape3.Move(d3)

		state = nextState
		stepNum++
	}

	return &SimulationResult{
		Steps:        steps,
		FinalState:   state,
		Accepted:     state == "q_accept",
		Tape1Initial: tape1Str,
		Tape2Initial: tape2Str,
		Tape3Initial: "",
	}
}
