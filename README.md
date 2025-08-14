# Modulo-Three FSM (Advanced Exercise)

## Overview
This project implements a deterministic finite-state machine (FSM) in Go to compute a binary number modulo 3.  
It is structured as a reusable FSM library plus a concrete “mod-3” example and a CLI tool for demonstration.

## Project Structure
```text
.
├── fsm/              # Generic FSM library (state mgmt, transitions, validation)
│   ├── fsm.go
│   └── fsm_test.go
├── modthree/         # Business logic: builds mod-3 FSM, exposes ModThree()
│   ├── modthree.go
│   └── modthree_test.go
└── cmd/modthree/     # CLI entry point using the modthree package
    ├── main.go
    └── main_test.go
```

## How It Works
1. FSM Library (fsm/)
   - Defines states, symbols, transitions, and execution flow (Run, Step, Process)
   - Validates machine setup before execution
   - Tracks current state and final/accepting states
2. Mod-3 Example (modthree/)
   - States: S0, S1, S2 represent remainder 0, 1, 2 modulo 3
   - Transitions follow binary mod-3 rules, MSB first
   - ModThree(input) trims whitespace, rejects non-binary characters, and returns the remainder
3. CLI (cmd/modthree/)
   - Accepts only `-in=<binary>` (no space between -in and value)
   - If value is quoted, leading/trailing spaces are trimmed

## Testing
1. Unit tests cover:
   - FSM library (fsm_test.go): validation errors, transition errors, symbol checks, empty inputs, accepting state checks
   - Mod-3 logic (modthree_test.go): correctness, whitespace handling, invalid input
   - CLI (main_test.go): valid/invalid args, help flag, spacing cases, error outputs
2. Run tests:
   ```bash
   go test ./... -cover
   ```
3. Results:
   ```text
   ok      assessment/cmd/modthree      coverage: 90.0% of statements
   ok      assessment/fsm               coverage: 95.9% of statements
   ok      assessment/modthree          coverage: 95.8% of statements
   ```

## Running the CLI
1. Help
   ```bash
   go run ./cmd/modthree -h
   ```
2. Valid input
   ```bash
   go run ./cmd/modthree -in=1101
   ```
3. Valid Output
   ```text
   modThree("1101") => 1
   ```