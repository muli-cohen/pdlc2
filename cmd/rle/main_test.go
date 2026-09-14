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

func runRLE(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
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

// TestCLIEncodeArg verifies AC3: positional argument encoding exits 0 with correct output.
func TestCLIEncodeArg(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "aaabbc")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "a3b2c1" {
		t.Errorf("got %q, want %q", got, "a3b2c1")
	}
}

// TestCLIDecodeArg verifies AC4: decode flag with positional argument exits 0 with correct output.
func TestCLIDecodeArg(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "-decode", "a3b2c1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "aaabbc" {
		t.Errorf("got %q, want %q", got, "aaabbc")
	}
}

// TestCLIMixedCase verifies AC5: mixed case is treated as distinct characters.
func TestCLIMixedCase(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "aAaa")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "a1A1a2" {
		t.Errorf("got %q, want %q", got, "a1A1a2")
	}
}

// TestCLIPunctuation verifies AC6: punctuation and whitespace are treated as ordinary characters.
func TestCLIPunctuation(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "hello, world!")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	want := "h1e1l2o1,1 1w1o1r1l1d1!1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestCLILongRun verifies AC7: 12 identical characters encode with a multi-digit count.
func TestCLILongRun(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", strings.Repeat("a", 12))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "a12" {
		t.Errorf("got %q, want %q", got, "a12")
	}
}

// TestCLIDecodeLongRun verifies AC8: decoding a multi-digit count expands to the correct number of characters.
func TestCLIDecodeLongRun(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "-decode", "a12")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	want := strings.Repeat("a", 12)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestCLIDigitInput verifies AC9: digit in encode input exits non-zero with a message mentioning digits.
func TestCLIDigitInput(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "abc123")
	if code == 0 {
		t.Fatalf("expected non-zero exit for digit input, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "digit") {
		t.Errorf("expected 'digit' in stderr, got: %s", stderr)
	}
}

// TestCLIMalformedDecode verifies AC10: malformed encoded input exits non-zero with a message about missing count.
func TestCLIMalformedDecode(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "-decode", "a3b")
	if code == 0 {
		t.Fatalf("expected non-zero exit for malformed input, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "count") {
		t.Errorf("expected 'count' in stderr, got: %s", stderr)
	}
}

// TestCLIStdin verifies AC11: stdin input is read and encoded correctly.
func TestCLIStdin(t *testing.T) {
	stdout, stderr, code := runRLE(t, "aaab\n")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "a3b1" {
		t.Errorf("got %q, want %q", got, "a3b1")
	}
}

// TestCLIConflict verifies AC12: both a positional argument and piped stdin exits non-zero with no output.
func TestCLIConflict(t *testing.T) {
	stdout, stderr, code := runRLE(t, "foo\n", "bar")
	if code == 0 {
		t.Fatalf("expected non-zero exit for arg+stdin conflict, got 0; stdout: %s", stdout)
	}
	if stdout != "" {
		t.Errorf("expected no stdout on conflict, got: %s", stdout)
	}
	if stderr == "" {
		t.Error("expected error message on stderr for conflict")
	}
}

// TestREADMEDocumentsUsage verifies AC14: README contains the required usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	examples := []string{
		`go run ./cmd/rle "aaabbc"`,
		`go run ./cmd/rle -decode "a3b2c1"`,
	}
	for _, ex := range examples {
		if !strings.Contains(content, ex) {
			t.Errorf("README.md missing example: %s", ex)
		}
	}
}
