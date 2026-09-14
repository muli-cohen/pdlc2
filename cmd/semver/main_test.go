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

func runSemver(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
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

func TestCompareMinor(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.2.3", "1.10.0")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "-1" {
		t.Errorf("got %q, want -1", strings.TrimSpace(stdout))
	}
}

func TestCompareBuildIgnored(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.0.0", "1.0.0+build.5")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "0" {
		t.Errorf("got %q, want 0", strings.TrimSpace(stdout))
	}
}

func TestComparePreReleaseVsRelease(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.0.0-alpha", "1.0.0")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "-1" {
		t.Errorf("got %q, want -1", strings.TrimSpace(stdout))
	}
}

func TestComparePreReleaseNumericVsAlpha(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.0.0-alpha.1", "1.0.0-alpha.beta")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "-1" {
		t.Errorf("got %q, want -1", strings.TrimSpace(stdout))
	}
}

func TestComparePreReleasePrefix(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.0.0-alpha", "1.0.0-alpha.1")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "-1" {
		t.Errorf("got %q, want -1", strings.TrimSpace(stdout))
	}
}

func TestSortMixed(t *testing.T) {
	stdout, stderr, code := runSemver(t, "1.10.0\n1.2.3\n1.0.0-rc.1\n1.0.0\n", "-sort")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	want := []string{"1.0.0-rc.1", "1.0.0", "1.2.3", "1.10.0"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d; output: %s", len(lines), len(want), stdout)
	}
	for i, line := range lines {
		if line != want[i] {
			t.Errorf("line %d: got %q, want %q", i, line, want[i])
		}
	}
}

func TestInvalidVersionStderr(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.01.0", "1.2.3")
	if code == 0 {
		t.Fatalf("expected non-zero exit, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "1.01.0") {
		t.Errorf("stderr must name the offending version, got: %s", stderr)
	}
}

func TestSortWithArgs(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "-sort", "1.0.0")
	if code == 0 {
		t.Fatalf("expected non-zero exit, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "-sort") {
		t.Errorf("stderr must mention -sort, got: %s", stderr)
	}
}

func TestWrongArgCount(t *testing.T) {
	stdout, stderr, code := runSemver(t, "", "1.0.0")
	if code == 0 {
		t.Fatalf("expected non-zero exit, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "expected 2") {
		t.Errorf("stderr must mention expected argument count, got: %s", stderr)
	}
}

func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	examples := []string{
		"go run ./cmd/semver 1.2.3 1.10.0",
		"printf '1.10.0\\n1.2.3\\n' | go run ./cmd/semver -sort",
	}
	for _, ex := range examples {
		if !strings.Contains(content, ex) {
			t.Errorf("README.md missing example: %s", ex)
		}
	}
}
