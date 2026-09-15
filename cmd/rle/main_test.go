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
	if strings.TrimSpace(stdout) != "a3b2c1" {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), "a3b2c1")
	}
}

// TestCLIDecodeArg verifies AC4: decode flag with positional argument exits 0 with correct output.
func TestCLIDecodeArg(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "-decode", "a3b2c1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "aaabbc" {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), "aaabbc")
	}
}

// TestCLIMixedCase verifies AC5: mixed case treated as distinct characters.
func TestCLIMixedCase(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "aAaa")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "a1A1a2" {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), "a1A1a2")
	}
}

// TestCLIPunctuation verifies AC6: whitespace and punctuation treated as ordinary characters.
func TestCLIPunctuation(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "hello, world!")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	want := "h1e1l2o1,1 1w1o1r1l1d1!1"
	if strings.TrimSpace(stdout) != want {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), want)
	}
}

// TestCLILongRun verifies AC7: runs longer than 9 characters.
func TestCLILongRun(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "aaaaaaaaaaaa")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "a12" {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), "a12")
	}
}

// TestCLIDecodeMultiDigit verifies AC8: decoding multi-digit counts.
func TestCLIDecodeMultiDigit(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "-decode", "a12")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "aaaaaaaaaaaa" {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), "aaaaaaaaaaaa")
	}
}

// TestCLIDigitError verifies AC9: digit in input exits non-zero with digit error.
func TestCLIDigitError(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "abc123")
	if code == 0 {
		t.Fatalf("expected non-zero exit for digit input, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "digit") {
		t.Errorf("expected \"digit\" in stderr, got: %s", stderr)
	}
}

// TestCLIMalformedMissingCount verifies AC10: malformed encoded string exits non-zero.
func TestCLIMalformedMissingCount(t *testing.T) {
	stdout, stderr, code := runRLE(t, "", "-decode", "a3b")
	if code == 0 {
		t.Fatalf("expected non-zero exit for malformed input, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "no following count") {
		t.Errorf("expected \"no following count\" in stderr, got: %s", stderr)
	}
}

// TestCLIStdin verifies AC11: stdin input with trailing newline stripped.
func TestCLIStdin(t *testing.T) {
	stdout, stderr, code := runRLE(t, "aaab\n")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "a3b1" {
		t.Errorf("got %q, want %q", strings.TrimSpace(stdout), "a3b1")
	}
}

// TestCLIConflict verifies AC12: positional argument and piped stdin exits non-zero.
func TestCLIConflict(t *testing.T) {
	stdout, stderr, code := runRLE(t, "aaab\n", "aaab")
	if code == 0 {
		t.Fatalf("expected non-zero exit for arg+stdin conflict, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "both") {
		t.Errorf("expected \"both\" in stderr, got: %s", stderr)
	}
}

// TestREADMEDocumentsUsage verifies AC15: README contains required usage examples and digit note.
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
	if !strings.Contains(content, "digit") {
		t.Error("README.md missing digit restriction note")
	}
}
