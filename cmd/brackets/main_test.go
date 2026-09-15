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
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func runBrackets(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// TestCLIBalancedArg verifies AC3: balanced input exits 0 and prints "balanced".
func TestCLIBalancedArg(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "", "(a[b]{c})")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "balanced" {
		t.Errorf("got %q, want %q", got, "balanced")
	}
}

// TestCLIUnbalancedArg verifies AC4: unbalanced input exits 0 and prints "unbalanced".
func TestCLIUnbalancedArg(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "", "(a[b)]")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "unbalanced" {
		t.Errorf("got %q, want %q", got, "unbalanced")
	}
}

// TestCLIPositionFlag verifies AC5: -position on a bad string prints the 0-based mismatch index.
func TestCLIPositionFlag(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "", "-position", "(a[b)]")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "4" {
		t.Errorf("got %q, want %q", got, "4")
	}
}

// TestCLINoBracketsDefault verifies AC6: no-bracket input in default mode prints "balanced".
func TestCLINoBracketsDefault(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "", "no brackets here")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "balanced" {
		t.Errorf("got %q, want %q", got, "balanced")
	}
}

// TestCLIPositionNoBrackets verifies AC7: -position on no-bracket input exits non-zero with error on stderr.
func TestCLIPositionNoBrackets(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "", "-position", "no brackets here")
	if code == 0 {
		t.Fatalf("expected non-zero exit, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "no brackets") {
		t.Errorf("expected error about no brackets on stderr, got: %s", stderr)
	}
}

// TestCLIStdin verifies AC8: input from stdin is read and balanced correctly.
func TestCLIStdin(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "{[()]}\n")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "balanced" {
		t.Errorf("got %q, want %q", got, "balanced")
	}
}

// TestCLIConflict verifies AC9: positional argument + piped stdin exits non-zero with error on stderr.
func TestCLIConflict(t *testing.T) {
	stdout, stderr, code := runBrackets(t, "hello\n", "world")
	if code == 0 {
		t.Fatalf("expected non-zero exit for arg+stdin conflict, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "not both") {
		t.Errorf("expected conflict error message in stderr, got: %s", stderr)
	}
}

// TestREADMEDocumentsUsage verifies AC13: README contains the three required usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	for _, example := range []string{
		`go run ./cmd/brackets "(a[b]{c})"`,
		`go run ./cmd/brackets -position "(a[b)]"`,
		`echo "{[()]}" | go run ./cmd/brackets`,
	} {
		if !strings.Contains(content, example) {
			t.Errorf("README.md missing example: %s", example)
		}
	}
}
