package modthree

import (
	"assessment/fsm"
	"fmt"
	"strings"
)

// States for the mod-3 FSM.
// S0, S1, S2 represent remainders 0, 1, and 2 when dividing by 3.
const (
	S0 fsm.State = "S0"
	S1 fsm.State = "S1"
	S2 fsm.State = "S2"
)

// BuildModThree creates an FSM that calculates a binary string's remainder mod 3.
// Input is processed MSB first, alphabet is {0,1}.
func BuildModThree() *fsm.FSM {
	f := fsm.NewFSM(S0)

	f.AddState(S0, true)
	f.AddState(S1, true)
	f.AddState(S2, true)

	f.AddSymbol('0')
	f.AddSymbol('1')

	// Transition table for remainders mod 3
	_ = f.AddTransition(S0, '0', S0)
	_ = f.AddTransition(S0, '1', S1)
	_ = f.AddTransition(S1, '0', S2)
	_ = f.AddTransition(S1, '1', S0)
	_ = f.AddTransition(S2, '0', S1)
	_ = f.AddTransition(S2, '1', S2)

	_ = f.Validate()
	return f
}

// ModThree returns the remainder (0, 1, or 2) of the binary string.
// Leading/trailing spaces are ignored; invalid symbols return an error.
func ModThree(input string) (int, error) {
	input = strings.TrimSpace(input)

	f := BuildModThree()
	final, err := f.Process(input)
	if err != nil {
		return -1, fmt.Errorf("invalid binary input: %w", err)
	}

	switch final {
	case S0:
		return 0, nil
	case S1:
		return 1, nil
	case S2:
		return 2, nil
	default:
		return -1, fmt.Errorf("unknown final state: %s", final)
	}
}
