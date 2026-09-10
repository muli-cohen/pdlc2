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

func runCaesar(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
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

// TestCLIEncodeArg verifies AC1/AC3: positional argument encoding exits 0 with correct output.
func TestCLIEncodeArg(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-shift", "3", "attack at dawn")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "dwwdfn dw gdzq" {
		t.Errorf("got %q, want %q", got, "dwwdfn dw gdzq")
	}
}

// TestCLIDecodeArg verifies AC2/AC4: decode flag with positional argument exits 0 with correct output.
func TestCLIDecodeArg(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-decode", "-shift", "3", "dwwdfn dw gdzq")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "attack at dawn" {
		t.Errorf("got %q, want %q", got, "attack at dawn")
	}
}

// TestCLIWrapAround verifies AC3/AC5: wrap-around from z to a and Z to A.
func TestCLIWrapAround(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-shift", "1", "xyz XYZ")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "yza YZA" {
		t.Errorf("got %q, want %q", got, "yza YZA")
	}
}

// TestCLIPassThrough verifies AC4/AC6: non-letter characters are unchanged.
func TestCLIPassThrough(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-shift", "3", "hello, world! 123")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "khoor, zruog! 123" {
		t.Errorf("got %q, want %q", got, "khoor, zruog! 123")
	}
}

// TestCLIShiftZero verifies AC5/AC7: shift 0 produces identity.
func TestCLIShiftZero(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-shift", "0", "abc")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

// TestCLIShift26 verifies shift 26 produces identity.
func TestCLIShift26(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-shift", "26", "abc")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

// TestCLIStdin verifies AC6/AC8: stdin input with default shift 3.
func TestCLIStdin(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "abc\n")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimSpace(stdout)
	if got != "def" {
		t.Errorf("got %q, want %q", got, "def")
	}
}

// TestCLIConflict verifies AC7/AC9: positional argument + piped stdin exits non-zero with error.
func TestCLIConflict(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "abc\n", "xyz")
	if code == 0 {
		t.Fatalf("expected non-zero exit for arg+stdin conflict, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "cannot accept both") {
		t.Errorf("expected conflict error message in stderr, got: %s", stderr)
	}
}

// TestCLIInvalidShift verifies AC10: a non-integer -shift value exits non-zero with a human-readable error.
func TestCLIInvalidShift(t *testing.T) {
	stdout, stderr, code := runCaesar(t, "", "-shift", "abc", "hello")
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid shift, got 0; stdout: %s", stdout)
	}
	if stderr == "" {
		t.Error("expected error message on stderr for invalid shift, got nothing")
	}
}

// TestREADMEDocumentsUsage verifies AC12: README contains the three required usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	examples := []string{
		`go run ./cmd/caesar -shift 3 "attack at dawn"`,
		`go run ./cmd/caesar -decode -shift 3 "dwwdfn dw gdzq"`,
		`echo "abc" | go run ./cmd/caesar`,
	}
	for _, ex := range examples {
		if !strings.Contains(content, ex) {
			t.Errorf("README.md missing example: %s", ex)
		}
	}
}
