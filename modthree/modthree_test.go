package modthree

import (
	"errors"
	"testing"

	"assessment/fsm"
)

// Test basic correct cases for ModThree.
func TestModThree_BasicCases(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"1101", 1},
		{"1110", 2},
		{"1111", 0},
		{"", 0}, // empty stays at S0
		{"0", 0},
		{"1", 1},
		{"10", 2},
		{"11", 0},
		{"01", 1}, // check MSB-first
		{"10", 2},
	}

	for _, tc := range tests {
		got, err := ModThree(tc.in)
		if err != nil {
			t.Fatalf("input=%q unexpected err: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("input=%q got=%d want=%d", tc.in, got, tc.want)
		}
	}
}

// Test trimming spaces in input.
func TestModThree_WhitespaceTrim(t *testing.T) {
	got, err := ModThree("   1110   ")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got != 2 {
		t.Fatalf("got=%d want=2", got)
	}
}

// Test invalid input returns fsm.ErrInvalidSymbol.
func TestModThree_InvalidInput(t *testing.T) {
	_, err := ModThree("10a01")
	if err == nil {
		t.Fatalf("expected error for non-binary input")
	}
	if !errors.Is(err, fsm.ErrInvalidSymbol) {
		t.Fatalf("want ErrInvalidSymbol, got %v", err)
	}
}

// Test the FSM from BuildModThree directly.
func TestBuildModThree_Sanity(t *testing.T) {
	m := BuildModThree()
	end, err := m.Process("1111") // should end in S0
	if err != nil {
		t.Fatalf("unexpected process error: %v", err)
	}
	if end != S0 {
		t.Fatalf("end=%s want=%s", end, S0)
	}
}
