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
	// working directory during go test is the package directory (cmd/joke)
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func TestCLIHappyPath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "run", ".")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("expected exit 0: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		t.Errorf("unexpected stderr output: %s", stderr.String())
	}
	joke := strings.TrimSpace(stdout.String())
	if joke == "" {
		t.Error("expected a non-empty joke on stdout, got empty string")
	}
}

func TestCLIExtraArgExitsNonZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "run", ".", "extra-arg")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit, got exit 0")
	}
	if stdout.Len() > 0 {
		t.Errorf("expected no stdout output, got: %s", stdout.String())
	}
	if stderr.Len() == 0 {
		t.Error("expected usage message on stderr, got nothing")
	}
}

func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	if !strings.Contains(string(data), "go run ./cmd/joke") {
		t.Error("README.md does not contain 'go run ./cmd/joke'")
	}
}
