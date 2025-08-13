package fsm

import (
	"errors"
	"testing"
)

//
// Validate()
//

func TestValidate_NoStates(t *testing.T) {
	// no states -> ErrNoStates
	f := NewFSM("S0")
	f.states = map[State]struct{}{} // simulate invalid build
	if err := f.Validate(); !errors.Is(err, ErrNoStates) {
		t.Fatalf("want ErrNoStates, got: %v", err)
	}
}

func TestValidate_EmptyAlphabet(t *testing.T) {
	// empty alphabet -> ErrEmptyAlphabet
	f := NewFSM("S0")
	f.AddState("S0", false)
	if err := f.Validate(); !errors.Is(err, ErrEmptyAlphabet) {
		t.Fatalf("want ErrEmptyAlphabet, got: %v", err)
	}
}

func TestValidate_StartNotInStates(t *testing.T) {
	// start not in states -> ErrStartNotInStates
	f := NewFSM("S0")
	f.AddState("S0", false)
	f.AddSymbol('0')
	f.start = "X"
	if err := f.Validate(); !errors.Is(err, ErrStartNotInStates) {
		t.Fatalf("want ErrStartNotInStates, got: %v", err)
	}
}

func TestValidate_OK(t *testing.T) {
	// valid config -> ok
	f := NewFSM("S0")
	f.AddState("S0", true)
	f.AddSymbol('0')
	if err := f.Validate(); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

//
// AddTransition()
//

func TestAddTransition_FromStateNotDefined(t *testing.T) {
	// unknown from-state -> error
	f := NewFSM("S")
	f.AddState("B", false)
	f.AddSymbol('0')
	if err := f.AddTransition("A", '0', "B"); err == nil {
		t.Fatal("expected error for from-state not defined")
	}
}

func TestAddTransition_ToStateNotDefined(t *testing.T) {
	// unknown to-state -> error
	f := NewFSM("A")
	f.AddState("A", false)
	f.AddSymbol('0')
	if err := f.AddTransition("A", '0', "B"); err == nil {
		t.Fatal("expected error for to-state not defined")
	}
}

func TestAddTransition_SymbolNotInAlphabet(t *testing.T) {
	// symbol not registered -> error
	f := NewFSM("A")
	f.AddState("A", false)
	f.AddState("B", false)
	if err := f.AddTransition("A", '1', "B"); err == nil {
		t.Fatalf("expected error for symbol not in alphabet")
	}
}

func TestAddTransition_Determinism(t *testing.T) {
	// duplicate (state,symbol) -> ErrDupTransition
	f := NewFSM("A")
	f.AddState("A", false)
	f.AddState("B", false)
	f.AddSymbol('0')

	if err := f.AddTransition("A", '0', "B"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := f.AddTransition("A", '0', "B"); !errors.Is(err, ErrDupTransition) {
		t.Fatalf("want ErrDupTransition, got: %v", err)
	}
}

//
// Step / Current / Accepting / Reset
//

func TestStep_InvalidSymbol(t *testing.T) {
	// invalid symbol -> ErrInvalidSymbol
	f := NewFSM("A")
	f.AddState("A", false)
	f.AddState("B", false)
	f.AddSymbol('0')
	_ = f.AddTransition("A", '0', "B")

	if err := f.Step('1', 0); !errors.Is(err, ErrInvalidSymbol) {
		t.Fatalf("want ErrInvalidSymbol, got: %v", err)
	}
}

func TestStep_NoTransition(t *testing.T) {
	// missing transition -> ErrNoTransition
	f := NewFSM("A")
	f.AddState("A", false)
	f.AddState("B", false)
	f.AddSymbol('0')

	if err := f.Step('0', 0); !errors.Is(err, ErrNoTransition) {
		t.Fatalf("want ErrNoTransition, got: %v", err)
	}
}

func TestRun_Current_Accepting_Reset(t *testing.T) {
	// A -'0'-> B ; B accepting
	f := NewFSM("A")
	f.AddState("A", false)
	f.AddState("B", true)
	f.AddSymbol('0')
	if err := f.AddTransition("A", '0', "B"); err != nil {
		t.Fatal(err)
	}

	// initial state
	if cur := f.Current(); cur != "A" {
		t.Fatalf("initial current=%s want=A", cur)
	}
	if f.Accepting() {
		t.Fatalf("start should not be accepting")
	}

	// run 1 step
	end, err := f.Run([]rune{'0'})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if end != "B" || f.Current() != "B" {
		t.Fatalf("end/current=%s want=B", end)
	}
	if !f.Accepting() {
		t.Fatalf("expected accepting at B")
	}

	// reset
	f.Reset()
	if cur := f.Current(); cur != "A" {
		t.Fatalf("after Reset current=%s want=A", cur)
	}
	if f.Accepting() {
		t.Fatalf("after Reset should not be accepting")
	}
}

//
// Process / Run flows
//

func TestProcess_Trim_RejectInternalSpace(t *testing.T) {
	// trims outer spaces; internal spaces invalid
	f := NewFSM("S")
	f.AddState("S", true)
	f.AddSymbol('0')
	if err := f.AddTransition("S", '0', "S"); err != nil {
		t.Fatal(err)
	}

	if _, err := f.Process("   0   "); err != nil {
		t.Fatalf("unexpected err on trimmed whitespace: %v", err)
	}
	if _, err := f.Process("0 0"); !errors.Is(err, ErrInvalidSymbol) {
		t.Fatalf("want ErrInvalidSymbol for internal space, got: %v", err)
	}
}

func TestProcess_EmptyInput_StayAtStart(t *testing.T) {
	// empty string stays at start
	f := NewFSM("S0")
	f.AddState("S0", true)
	f.AddSymbol('0')
	end, err := f.Process("")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if end != "S0" {
		t.Fatalf("end=%s want=S0", end)
	}
}

func TestRun_EmptySeq_StayAtStart(t *testing.T) {
	// empty rune slice stays at start
	f := NewFSM("S0")
	f.AddState("S0", false)
	f.AddSymbol('0')

	end, err := f.Run(nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if end != "S0" || f.Current() != "S0" {
		t.Fatalf("end/current=%s want=S0", end)
	}
	if f.Accepting() {
		t.Fatalf("should not be accepting at start")
	}
}

//
// IsFinal
//

func TestIsFinal_TrueFalse(t *testing.T) {
	// true for final; false otherwise
	f := NewFSM("X")
	f.AddState("A", false)
	f.AddState("B", true)
	if f.IsFinal("A") {
		t.Fatal("A should not be final")
	}
	if !f.IsFinal("B") {
		t.Fatal("B should be final")
	}
}
