package fsm

import (
	"errors"
	"fmt"
	"strings"
)

// State is a named state.
type State string

// Symbol is a single input symbol.
type Symbol rune

// Common errors.
var (
	ErrStartNotInStates = errors.New("start state not in states set")
	ErrInvalidSymbol    = errors.New("invalid symbol")
	ErrNoTransition     = errors.New("no transition")
	ErrDupTransition    = errors.New("transition already defined")
	ErrNoStates         = errors.New("no states defined")
	ErrEmptyAlphabet    = errors.New("alphabet is empty")
)

// FSM is a deterministic finite-state machine.
type FSM struct {
	states   map[State]struct{}
	alphabet map[Symbol]struct{}
	start    State
	finals   map[State]struct{}
	trans    map[State]map[Symbol]State
	current  State
}

// NewFSM creates an FSM with the given start state.
func NewFSM(start State) *FSM {
	return &FSM{
		states:   map[State]struct{}{start: {}},
		alphabet: make(map[Symbol]struct{}),
		start:    start,
		finals:   make(map[State]struct{}),
		trans:    make(map[State]map[Symbol]State),
		current:  start,
	}
}

// AddState adds a state; mark as final if isFinal is true.
func (f *FSM) AddState(s State, isFinal bool) {
	f.states[s] = struct{}{}
	if isFinal {
		f.finals[s] = struct{}{}
	}
	if _, ok := f.trans[s]; !ok {
		f.trans[s] = make(map[Symbol]State)
	}
}

// AddSymbol registers a symbol in the alphabet.
func (f *FSM) AddSymbol(r rune) {
	f.alphabet[Symbol(r)] = struct{}{}
}

// AddTransition sets δ(from, symbol) → to; errors if invalid or duplicate.
func (f *FSM) AddTransition(from State, symbol rune, to State) error {
	if _, ok := f.states[from]; !ok {
		return fmt.Errorf("from-state %q not defined", from)
	}
	if _, ok := f.states[to]; !ok {
		return fmt.Errorf("to-state %q not defined", to)
	}
	s := Symbol(symbol)
	if _, ok := f.alphabet[s]; !ok {
		return fmt.Errorf("symbol %q not in alphabet", rune(symbol))
	}
	if _, ok := f.trans[from]; !ok {
		f.trans[from] = make(map[Symbol]State)
	}
	if _, exists := f.trans[from][s]; exists {
		return fmt.Errorf("%w for (%s,%q)", ErrDupTransition, from, symbol)
	}
	f.trans[from][s] = to
	return nil
}

// Validate ensures states/alphabet exist and start is valid.
func (f *FSM) Validate() error {
	if len(f.states) == 0 {
		return ErrNoStates
	}
	if len(f.alphabet) == 0 {
		return ErrEmptyAlphabet
	}
	if _, ok := f.states[f.start]; !ok {
		return ErrStartNotInStates
	}
	return nil
}

// IsFinal reports whether s is final.
func (f *FSM) IsFinal(s State) bool {
	_, ok := f.finals[s]
	return ok
}

// Reset sets current to start.
func (f *FSM) Reset() { f.current = f.start }

// Step consumes one symbol and moves to next state.
func (f *FSM) Step(r rune, pos int) error {
	s := Symbol(r)
	if _, ok := f.alphabet[s]; !ok {
		return fmt.Errorf("%w: %q at position %d", ErrInvalidSymbol, r, pos)
	}
	next, ok := f.trans[f.current][s]
	if !ok {
		return fmt.Errorf("%w: from %s on %q at pos %d", ErrNoTransition, f.current, r, pos)
	}
	f.current = next
	return nil
}

// Current returns the current state.
func (f *FSM) Current() State { return f.current }

// Accepting reports whether current state is final.
func (f *FSM) Accepting() bool { return f.IsFinal(f.current) }

// Run executes the FSM over the rune sequence in order.
func (f *FSM) Run(seq []rune) (State, error) {
	if err := f.Validate(); err != nil {
		return "", err
	}
	f.Reset()
	for i, r := range seq {
		if err := f.Step(r, i); err != nil {
			return "", err
		}
	}
	return f.current, nil
}

// Process executes the FSM over a string (in order).
// Trims leading/trailing spaces; internal spaces are invalid.
func (f *FSM) Process(input string) (State, error) {
	input = strings.TrimSpace(input)
	return f.Run([]rune(input))
}
