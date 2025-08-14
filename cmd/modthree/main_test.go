package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"assessment/modthree"
)

// runWith executes run() with given args and returns captured stdout, stderr, and exit code.
func runWith(args ...string) (stdout, stderr string, exit int) {
	var out, err bytes.Buffer
	code := run(args, &out, &err)
	return out.String(), err.String(), code
}

// containsAll reports whether all substrings are present in the given string.
func containsAll(haystack string, needles ...string) bool {
	for _, n := range needles {
		if !strings.Contains(haystack, n) {
			return false
		}
	}
	return true
}

func TestRun_Valid_NoSpaces(t *testing.T) {
	stdout, stderr, code := runWith("-in=1101")
	if code != ExitOK {
		t.Fatalf("exit=%d, want %d, stderr=%q", code, ExitOK, stderr)
	}
	if stderr != "" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
	wantRes, err := modthree.ModThree("1101")
	if err != nil {
		t.Fatalf("ModThree error: %v", err)
	}
	if !containsAll(stdout, `modThree("1101") => `, strconv.Itoa(wantRes)) {
		t.Fatalf("stdout=%q, want contains modThree(\"1101\") => %d", stdout, wantRes)
	}
}

func TestRun_Valid_TrimSpacesInQuotedValue(t *testing.T) {
	stdout, stderr, code := runWith("-in=   1101   ")
	if code != ExitOK {
		t.Fatalf("exit=%d, want %d, stderr=%q", code, ExitOK, stderr)
	}
	if stderr != "" {
		t.Fatalf("unexpected stderr: %q", stderr)
	}
	wantRes, _ := modthree.ModThree("1101")
	if !containsAll(stdout, `modThree("1101") => `, strconv.Itoa(wantRes)) {
		t.Fatalf("stdout=%q, want contains modThree(\"1101\") => %d", stdout, wantRes)
	}
}

func TestRun_Invalid_SpaceSeparated(t *testing.T) {
	_, stderr, code := runWith("-in", "1101")
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !(strings.Contains(stderr, "invalid format") || strings.Contains(stderr, "use -in=<binary>")) {
		t.Fatalf("stderr=%q, want message about invalid format", stderr)
	}
}

func TestRun_Invalid_EqualsThenSpace_NoQuotes(t *testing.T) {
	_, stderr, code := runWith("-in=", "1101")
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !(strings.Contains(stderr, "unexpected extra args") || strings.Contains(stderr, "provide -in=<binary>")) {
		t.Fatalf("stderr=%q, want message about extra args or missing value", stderr)
	}
}

func TestRun_Invalid_NonBinary(t *testing.T) {
	_, stderr, code := runWith("-in=1102")
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !strings.Contains(stderr, "must contain only 0 and 1") {
		t.Fatalf("stderr=%q, want validation message", stderr)
	}
}

func TestRun_Invalid_MissingValueFlagPresent(t *testing.T) {
	_, stderr, code := runWith("-in=")
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !strings.Contains(stderr, "provide -in=<binary>") {
		t.Fatalf("stderr=%q, want prompt to provide -in", stderr)
	}
}

func TestRun_Invalid_MissingFlagCompletely(t *testing.T) {
	_, stderr, code := runWith()
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !strings.Contains(stderr, "provide -in=<binary>") {
		t.Fatalf("stderr=%q, want prompt to provide -in", stderr)
	}
}

func TestRun_Invalid_ExtraPositionalArgs(t *testing.T) {
	_, stderr, code := runWith("-in=101", "extra")
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !strings.Contains(stderr, "unexpected extra args") {
		t.Fatalf("stderr=%q, want message about unexpected extra args", stderr)
	}
}

func TestRun_ParseError_UnknownFlag(t *testing.T) {
	_, stderr, code := runWith("-x")
	if code != ExitInvalidArgs {
		t.Fatalf("exit=%d, want %d", code, ExitInvalidArgs)
	}
	if !(strings.Contains(stderr, "flag provided but not defined") || strings.Contains(stderr, "Usage: modthree -in=<binary>")) {
		t.Fatalf("stderr=%q, want flag parse error output", stderr)
	}
}

func TestRun_Help_IsSuccess(t *testing.T) {
	_, stderr, code := runWith("-h")
	if code != ExitOK {
		t.Fatalf("exit=%d, want %d", code, ExitOK)
	}
	if !strings.Contains(stderr, "Usage: modthree -in=<binary>") {
		t.Fatalf("stderr=%q, want usage output", stderr)
	}

	_, stderr, code = runWith("-help")
	if code != ExitOK {
		t.Fatalf("exit=%d, want %d", code, ExitOK)
	}
	if !strings.Contains(stderr, "Usage: modthree -in=<binary>") {
		t.Fatalf("stderr=%q, want usage output", stderr)
	}
}
