package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// symbolChars mirrors the symbol alphabet defined in the password package.
const symbolChars = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func runPassgen(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
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

// TestCLIDefaultOutput verifies AC3: default run prints a 16-character password to stdout and exits 0.
func TestCLIDefaultOutput(t *testing.T) {
	stdout, stderr, code := runPassgen(t)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("expected no stderr output, got: %s", stderr)
	}
	pw := strings.TrimSpace(stdout)
	if len(pw) != 16 {
		t.Errorf("expected 16-character password, got %d characters: %q", len(pw), pw)
	}
}

// TestCLICustomLength verifies AC4: -length 20 prints a 20-character password and exits 0.
func TestCLICustomLength(t *testing.T) {
	stdout, stderr, code := runPassgen(t, "-length", "20")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	pw := strings.TrimSpace(stdout)
	if len(pw) != 20 {
		t.Errorf("expected 20-character password, got %d characters: %q", len(pw), pw)
	}
}

// TestCLINoSymbols verifies AC5: -symbols=false produces a password with no symbol characters and exits 0.
func TestCLINoSymbols(t *testing.T) {
	stdout, stderr, code := runPassgen(t, "-symbols=false")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	pw := strings.TrimSpace(stdout)
	if strings.ContainsAny(pw, symbolChars) {
		t.Errorf("expected no symbol characters in %q, but found some", pw)
	}
}

// TestCLILengthBelowMinimum verifies AC6: -length 4 exits non-zero, stderr mentions minimum length, stdout is empty.
func TestCLILengthBelowMinimum(t *testing.T) {
	stdout, stderr, code := runPassgen(t, "-length", "4")
	if code == 0 {
		t.Fatal("expected non-zero exit code for -length 4, got 0")
	}
	if stdout != "" {
		t.Errorf("expected nothing on stdout, got: %s", stdout)
	}
	if !strings.Contains(stderr, "8") {
		t.Errorf("expected stderr to mention minimum length 8, got: %s", stderr)
	}
}

// TestREADMEDocumentsUsage verifies AC11: README contains passgen usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	for _, example := range []string{
		"go run ./cmd/passgen",
		"go run ./cmd/passgen -length 20",
		"go run ./cmd/passgen -symbols=false",
		"go run ./cmd/passgen -length 4",
	} {
		if !strings.Contains(content, example) {
			t.Errorf("README.md missing example: %s", example)
		}
	}
}
