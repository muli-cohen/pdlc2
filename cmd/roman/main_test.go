package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	// working directory during go test is the package directory (cmd/roman)
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func runRoman(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// TestCLIIntegerToNumeral verifies AC3: integer argument converts to numeral on stdout, exits 0.
func TestCLIIntegerToNumeral(t *testing.T) {
	stdout, stderr, code := runRoman(t, "1994")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("expected no stderr output, got: %s", stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "MCMXCIV" {
		t.Errorf("expected %q, got %q", "MCMXCIV", got)
	}
}

// TestCLINumeralToInteger verifies AC4: numeral argument converts to integer on stdout, exits 0.
func TestCLINumeralToInteger(t *testing.T) {
	stdout, stderr, code := runRoman(t, "MCMXCIV")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("expected no stderr output, got: %s", stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "1994" {
		t.Errorf("expected %q, got %q", "1994", got)
	}
}

// TestCLICaseInsensitive verifies AC5: lowercase numeral is accepted and converts correctly.
func TestCLICaseInsensitive(t *testing.T) {
	stdout, stderr, code := runRoman(t, "mcmxciv")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("expected no stderr output, got: %s", stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "1994" {
		t.Errorf("expected %q, got %q", "1994", got)
	}
}

// TestCLIOutOfRangeLow verifies AC6: integer 0 exits non-zero; output mentions the supported range.
func TestCLIOutOfRangeLow(t *testing.T) {
	stdout, stderr, code := runRoman(t, "0")
	if code == 0 {
		t.Fatal("expected non-zero exit for 0, got 0")
	}
	if stdout != "" {
		t.Errorf("expected no stdout output, got: %s", stdout)
	}
	if !strings.Contains(stderr, "1-3999") {
		t.Errorf("expected stderr to mention supported range 1-3999, got: %s", stderr)
	}
}

// TestCLIOutOfRangeHigh verifies AC7: integer 4000 exits non-zero; output mentions the supported range.
func TestCLIOutOfRangeHigh(t *testing.T) {
	stdout, stderr, code := runRoman(t, "4000")
	if code == 0 {
		t.Fatal("expected non-zero exit for 4000, got 0")
	}
	if stdout != "" {
		t.Errorf("expected no stdout output, got: %s", stdout)
	}
	if !strings.Contains(stderr, "1-3999") {
		t.Errorf("expected stderr to mention supported range 1-3999, got: %s", stderr)
	}
}

// TestCLINonCanonicalNumeral verifies AC8: non-canonical numeral exits non-zero with an error message.
func TestCLINonCanonicalNumeral(t *testing.T) {
	stdout, stderr, code := runRoman(t, "IIII")
	if code == 0 {
		t.Fatal("expected non-zero exit for IIII, got 0")
	}
	if stdout != "" {
		t.Errorf("expected no stdout output, got: %s", stdout)
	}
	if stderr == "" {
		t.Error("expected error message on stderr, got nothing")
	}
}

// TestCLINoArgs verifies usage message and non-zero exit when called with no arguments.
func TestCLINoArgs(t *testing.T) {
	stdout, stderr, code := runRoman(t)
	if code == 0 {
		t.Fatal("expected non-zero exit with no arguments, got 0")
	}
	if stdout != "" {
		t.Errorf("expected no stdout output, got: %s", stdout)
	}
	if stderr == "" {
		t.Error("expected usage message on stderr, got nothing")
	}
}

// TestCLITooManyArgs verifies usage message and non-zero exit when called with more than one argument.
func TestCLITooManyArgs(t *testing.T) {
	stdout, stderr, code := runRoman(t, "1994", "extra")
	if code == 0 {
		t.Fatal("expected non-zero exit with two arguments, got 0")
	}
	if stdout != "" {
		t.Errorf("expected no stdout output, got: %s", stdout)
	}
	if stderr == "" {
		t.Error("expected usage message on stderr, got nothing")
	}
}

// TestREADMEDocumentsUsage verifies AC13: README contains the required usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	for _, example := range []string{
		"go run ./cmd/roman 1994",
		"go run ./cmd/roman MCMXCIV",
	} {
		if !strings.Contains(content, example) {
			t.Errorf("README.md missing example: %s", example)
		}
	}
}
