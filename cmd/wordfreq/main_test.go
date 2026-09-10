package main_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"os/exec"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func runWordfreq(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Stdin = strings.NewReader(stdin)
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

// TestCLIBasicStdin verifies AC3: pipe "the cat the hat", check stdout lines.
func TestCLIBasicStdin(t *testing.T) {
	stdout, stderr, code := runWordfreq(t, "the cat the hat")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	want := []string{"2 the", "1 cat", "1 hat"}
	if len(lines) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(lines), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i, lines[i], w)
		}
	}
}

// TestCLIFlagN verifies AC4: -n 1 on the same input prints only "2 the".
func TestCLIFlagN(t *testing.T) {
	stdout, stderr, code := runWordfreq(t, "the cat the hat", "-n", "1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != "2 the" {
		t.Errorf("expected %q, got %q", "2 the", got)
	}
}

// TestCLIEmptyInput verifies AC5: empty input produces no stdout and exits 0.
func TestCLIEmptyInput(t *testing.T) {
	stdout, stderr, code := runWordfreq(t, "")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stdout != "" {
		t.Errorf("expected empty stdout, got: %q", stdout)
	}
}

// TestCLIMissingFile verifies AC6: missing file exits non-zero and stderr names the file.
func TestCLIMissingFile(t *testing.T) {
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("go", "run", ".", "missing.txt")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing file, got 0")
	}
	stderr := errBuf.String()
	if !strings.Contains(stderr, "missing.txt") {
		t.Errorf("expected stderr to contain %q, got: %s", "missing.txt", stderr)
	}
}

// TestCLIInvalidN verifies AC7: -n 0 exits non-zero and stderr mentions -n requirement.
func TestCLIInvalidN(t *testing.T) {
	stdout, stderr, code := runWordfreq(t, "the cat", "-n", "0")
	if code == 0 {
		t.Fatal("expected non-zero exit for -n 0, got 0")
	}
	if stdout != "" {
		t.Errorf("expected empty stdout, got: %q", stdout)
	}
	if !strings.Contains(stderr, "-n") {
		t.Errorf("expected stderr to mention -n, got: %s", stderr)
	}
}

// TestREADMEDocumentsUsage verifies AC9: README contains required usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	for _, example := range []string{
		"go run ./cmd/wordfreq README.md",
		"go run ./cmd/wordfreq -n 5",
	} {
		if !strings.Contains(content, example) {
			t.Errorf("README.md missing example: %s", example)
		}
	}
}
