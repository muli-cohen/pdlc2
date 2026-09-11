package main_test

import (
	"bytes"
	"encoding/json"
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

func runSubnet(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
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

// TestCLITextOutput verifies AC3: correct ordered text output for 192.168.1.10/24.
func TestCLITextOutput(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "192.168.1.10/24")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	want := []string{
		"Address: 192.168.1.10",
		"Network: 192.168.1.0",
		"Broadcast: 192.168.1.255",
		"Mask: 255.255.255.0",
		"FirstHost: 192.168.1.1",
		"LastHost: 192.168.1.254",
		"PrefixLength: 24",
		"UsableHosts: 254",
	}
	if len(lines) != len(want) {
		t.Fatalf("expected %d lines, got %d:\n%s", len(want), len(lines), stdout)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d: got %q, want %q", i+1, lines[i], w)
		}
	}
}

// TestCLIHostBitMasking verifies AC4: 192.168.1.200/24 produces the same Network as 192.168.1.10/24.
func TestCLIHostBitMasking(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "192.168.1.200/24")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "Network: 192.168.1.0") {
		t.Errorf("expected Network: 192.168.1.0 in output:\n%s", stdout)
	}
}

// TestCLISlash32 verifies AC5: /32 omits Broadcast, FirstHost, LastHost lines and reports 0 usable hosts.
func TestCLISlash32(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "10.0.0.1/32")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	for _, absent := range []string{"Broadcast:", "FirstHost:", "LastHost:"} {
		if strings.Contains(stdout, absent) {
			t.Errorf("expected %s line to be absent for /32, but found it in output:\n%s", absent, stdout)
		}
	}
	if !strings.Contains(stdout, "UsableHosts: 0") {
		t.Errorf("expected UsableHosts: 0 in output:\n%s", stdout)
	}
}

// TestCLISlash8UsableHosts verifies AC6: /8 reports 16777214 usable hosts.
func TestCLISlash8UsableHosts(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "10.0.0.0/8")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "UsableHosts: 16777214") {
		t.Errorf("expected UsableHosts: 16777214 in output:\n%s", stdout)
	}
}

// TestCLIJSONOutput verifies AC7: -json produces valid JSON with correct fields for /8.
func TestCLIJSONOutput(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "-json", "10.0.0.0/8")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, stdout)
	}
	checks := map[string]interface{}{
		"Address":      "10.0.0.0",
		"Network":      "10.0.0.0",
		"Broadcast":    "10.255.255.255",
		"Mask":         "255.0.0.0",
		"FirstHost":    "10.0.0.1",
		"LastHost":     "10.255.255.254",
		"PrefixLength": float64(8),
		"UsableHosts":  float64(16777214),
	}
	for k, want := range checks {
		got, ok := obj[k]
		if !ok {
			t.Errorf("JSON missing key %q", k)
			continue
		}
		if got != want {
			t.Errorf("JSON[%q]: got %v, want %v", k, got, want)
		}
	}
}

// TestCLIJSONSlash32OmitsAbsentFields verifies /32 JSON omits Broadcast, FirstHost, LastHost.
func TestCLIJSONSlash32OmitsAbsentFields(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "-json", "10.0.0.1/32")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d; stderr: %s", code, stderr)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, stdout)
	}
	for _, absent := range []string{"Broadcast", "FirstHost", "LastHost"} {
		if _, ok := obj[absent]; ok {
			t.Errorf("JSON key %q should be absent for /32, but it is present", absent)
		}
	}
}

// TestCLIInvalidPrefixExitsNonZero verifies AC8: /33 exits non-zero with a message naming the offending value.
func TestCLIInvalidPrefixExitsNonZero(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "192.168.1.1/33")
	if code == 0 {
		t.Fatalf("expected non-zero exit for /33, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "33") {
		t.Errorf("expected stderr to mention '33', got: %s", stderr)
	}
}

// TestCLIInvalidOctetExitsNonZero verifies AC9: octet 300 exits non-zero with a message naming the offending value.
func TestCLIInvalidOctetExitsNonZero(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "300.1.1.1/24")
	if code == 0 {
		t.Fatalf("expected non-zero exit for octet 300, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "300") {
		t.Errorf("expected stderr to mention '300', got: %s", stderr)
	}
}

// TestCLINoArgsExitsNonZero verifies AC10: no arguments exits non-zero with a usage hint on stderr.
func TestCLINoArgsExitsNonZero(t *testing.T) {
	stdout, stderr, code := runSubnet(t)
	if code == 0 {
		t.Fatalf("expected non-zero exit for missing arg, got 0; stdout: %s", stdout)
	}
	if stderr == "" {
		t.Errorf("expected usage message on stderr, got nothing")
	}
}

// TestCLIExtraArgExitsNonZero verifies AC11: extra positional argument exits non-zero naming the offending arg.
func TestCLIExtraArgExitsNonZero(t *testing.T) {
	stdout, stderr, code := runSubnet(t, "192.168.1.1/24", "10.0.0.1/8")
	if code == 0 {
		t.Fatalf("expected non-zero exit for extra arg, got 0; stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "10.0.0.1/8") {
		t.Errorf("expected stderr to name extra argument '10.0.0.1/8', got: %s", stderr)
	}
}

// TestREADMEDocumentsUsage verifies AC12: README contains the two required usage examples.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	examples := []string{
		"go run ./cmd/subnet 192.168.1.10/24",
		"go run ./cmd/subnet -json 10.0.0.0/8",
	}
	for _, ex := range examples {
		if !strings.Contains(content, ex) {
			t.Errorf("README.md missing example: %s", ex)
		}
	}
}
