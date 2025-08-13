package main

import (
	"assessment/modthree"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Exit codes
const (
	ExitOK          = 0 // success
	ExitRuntimeErr  = 1 // computation error
	ExitInvalidArgs = 2 // invalid CLI arguments
)

// isBinary returns true if s is non-empty and contains only '0'/'1'.
func isBinary(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r != '0' && r != '1' {
			return false
		}
	}
	return true
}

// run parses CLI args, validates them, runs the computation, and returns an exit code.
// Accepts only "-in=<value>" form; quoted values are trimmed.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("modthree", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "Usage: modthree -in=<binary>")
		_, _ = fmt.Fprintln(stderr, "  Accepts only -in=<binary> (no space between -in and value)")
		_, _ = fmt.Fprintln(stderr, "  If value is quoted, leading/trailing spaces are trimmed")
		_, _ = fmt.Fprintln(stderr, "  Example: modthree -in=1101  or  modthree -in=\"   1101   \"")
	}

	// Reject "-in <value>" (space-separated)
	for _, a := range args {
		if a == "-in" || strings.HasPrefix(a, "-in ") {
			_, _ = fmt.Fprintln(stderr, "invalid format: use -in=<binary> with no space")
			return ExitInvalidArgs
		}
	}

	var input string
	fs.StringVar(&input, "in", "", "binary input string (MSB first), e.g. 1101")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitOK
		}
		return ExitInvalidArgs
	}

	if fs.NArg() > 0 {
		_, _ = fmt.Fprintf(stderr, "unexpected extra args: %v\n", fs.Args())
		return ExitInvalidArgs
	}

	input = strings.TrimSpace(input)
	if input == "" {
		_, _ = fmt.Fprintln(stderr, "provide -in=<binary>")
		return ExitInvalidArgs
	}
	if !isBinary(input) {
		_, _ = fmt.Fprintf(stderr, "invalid input %q: must contain only 0 and 1\n", input)
		return ExitInvalidArgs
	}

	res, err := modthree.ModThree(input)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "error:", err)
		return ExitRuntimeErr
	}

	_, _ = fmt.Fprintf(stdout, "modThree(%q) => %d\n", input, res)
	return ExitOK
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
