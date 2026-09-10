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
	// working directory during go test is the package directory (cmd/todo)
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

// run invokes the cmd/todo binary via go run . with the given -file path and subcommand args.
// file must be an absolute path; the -file flag is always passed so tests don't litter the
// package directory with todo.json.
func run(t *testing.T, file string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	goArgs := []string{"run", "."}
	if file != "" {
		goArgs = append(goArgs, "-file", file)
	}
	goArgs = append(goArgs, args...)
	cmd := exec.Command("go", goArgs...)
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

// TestCLIAddAndList verifies AC3 and AC4: add creates the file, prints the ID;
// list shows the pending item in the correct format.
func TestCLIAddAndList(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	stdout, stderr, code := run(t, file, "add", "buy milk")
	if code != 0 {
		t.Fatalf("add: expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stderr != "" {
		t.Errorf("add: unexpected stderr: %s", stderr)
	}
	if !strings.Contains(stdout, "1") {
		t.Errorf("add: expected item ID 1 in output, got: %s", stdout)
	}
	if _, err := os.Stat(file); err != nil {
		t.Errorf("add: expected JSON file at %s, got: %v", file, err)
	}

	stdout, stderr, code = run(t, file, "list")
	if code != 0 {
		t.Fatalf("list: expected exit 0, got %d; stderr: %s", code, stderr)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != "1 [ ] buy milk" {
		t.Errorf("list: expected %q, got %q", "1 [ ] buy milk", got)
	}
}

// TestCLIDone verifies AC5: done marks the item; list then shows [x].
func TestCLIDone(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	run(t, file, "add", "buy milk")

	stdout, stderr, code := run(t, file, "done", "1")
	if code != 0 {
		t.Fatalf("done: expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stdout != "" {
		t.Errorf("done: expected no stdout, got: %s", stdout)
	}

	stdout, _, code = run(t, file, "list")
	if code != 0 {
		t.Fatalf("list after done: expected exit 0, got %d", code)
	}
	got := strings.TrimRight(stdout, "\n")
	if got != "1 [x] buy milk" {
		t.Errorf("list after done: expected %q, got %q", "1 [x] buy milk", got)
	}
}

// TestCLIRemove verifies AC6: rm removes the item; list then prints nothing and exits 0.
func TestCLIRemove(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	run(t, file, "add", "buy milk")

	stdout, stderr, code := run(t, file, "rm", "1")
	if code != 0 {
		t.Fatalf("rm: expected exit 0, got %d; stderr: %s", code, stderr)
	}
	if stdout != "" {
		t.Errorf("rm: expected no stdout, got: %s", stdout)
	}

	stdout, _, code = run(t, file, "list")
	if code != 0 {
		t.Fatalf("list after rm: expected exit 0, got %d", code)
	}
	if stdout != "" {
		t.Errorf("list after rm: expected no output, got: %s", stdout)
	}
}

// TestCLIDoneUnknownID verifies AC7: done with an unknown ID exits non-zero and names the ID.
func TestCLIDoneUnknownID(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	stdout, stderr, code := run(t, file, "done", "99")
	if code == 0 {
		t.Fatal("done 99: expected non-zero exit, got 0")
	}
	if stdout != "" {
		t.Errorf("done 99: expected no stdout, got: %s", stdout)
	}
	if !strings.Contains(stderr, "99") {
		t.Errorf("done 99: expected stderr to name id 99, got: %s", stderr)
	}
}

// TestCLIAddEmptyTitle verifies AC8: add "" exits non-zero with an error: prefix on stderr.
func TestCLIAddEmptyTitle(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	stdout, stderr, code := run(t, file, "add", "")
	if code == 0 {
		t.Fatal(`add "": expected non-zero exit, got 0`)
	}
	if stdout != "" {
		t.Errorf(`add "": expected no stdout, got: %s`, stdout)
	}
	if !strings.HasPrefix(stderr, "error:") {
		t.Errorf(`add "": expected stderr to start with "error:", got: %s`, stderr)
	}
}

// TestCLIUnknownSubcommand verifies AC9: an unknown subcommand exits non-zero with usage on stderr.
func TestCLIUnknownSubcommand(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	stdout, stderr, code := run(t, file, "foo")
	if code == 0 {
		t.Fatal("foo: expected non-zero exit, got 0")
	}
	if stdout != "" {
		t.Errorf("foo: expected no stdout, got: %s", stdout)
	}
	if stderr == "" {
		t.Error("foo: expected usage message on stderr, got nothing")
	}
}

// TestCLIMissingArgument verifies AC10: add, done, and rm called with no argument
// each exit non-zero and print a usage message to stderr.
func TestCLIMissingArgument(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	for _, sub := range []string{"add", "done", "rm"} {
		stdout, stderr, code := run(t, file, sub)
		if code == 0 {
			t.Errorf("%s (no arg): expected non-zero exit, got 0", sub)
		}
		if stdout != "" {
			t.Errorf("%s (no arg): expected no stdout, got: %s", sub, stdout)
		}
		if stderr == "" {
			t.Errorf("%s (no arg): expected usage on stderr, got nothing", sub)
		}
	}
}

// TestCLIIDMonotonicity verifies AC11: after add A, add B, rm A, add C; list shows IDs 2 and 3.
func TestCLIIDMonotonicity(t *testing.T) {
	file := filepath.Join(t.TempDir(), "todo.json")

	run(t, file, "add", "A")
	run(t, file, "add", "B")
	run(t, file, "rm", "1")
	run(t, file, "add", "C")

	stdout, stderr, code := run(t, file, "list")
	if code != 0 {
		t.Fatalf("list: expected exit 0, got %d; stderr: %s", code, stderr)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 items, got %d: %v", len(lines), lines)
	}
	if !strings.HasPrefix(lines[0], "2 ") {
		t.Errorf("expected first item ID=2, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "3 ") {
		t.Errorf("expected second item ID=3, got: %s", lines[1])
	}
}

// TestCLIFileFlag verifies AC12: -file redirects storage to the given path,
// keeping separate lists fully isolated from one another.
func TestCLIFileFlag(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "list1.json")
	file2 := filepath.Join(dir, "list2.json")

	run(t, file1, "add", "item in list one")
	run(t, file2, "add", "item in list two")

	// list1 must contain only its own item.
	stdout, _, code := run(t, file1, "list")
	if code != 0 {
		t.Fatalf("list file1: expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "list one") || strings.Contains(stdout, "list two") {
		t.Errorf("list file1: expected only 'list one' item, got: %s", stdout)
	}

	// list2 must contain only its own item.
	stdout, _, code = run(t, file2, "list")
	if code != 0 {
		t.Fatalf("list file2: expected exit 0, got %d", code)
	}
	if !strings.Contains(stdout, "list two") || strings.Contains(stdout, "list one") {
		t.Errorf("list file2: expected only 'list two' item, got: %s", stdout)
	}

	// The default todo.json must not have been created in the package working directory
	// (all invocations above used an explicit -file flag).
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	if _, err := os.Stat(filepath.Join(wd, "todo.json")); err == nil {
		t.Error("default todo.json was created in the working directory; -file flag was not respected")
	}
}

// TestREADMEDocumentsUsage verifies AC13: README contains runnable examples for all subcommands
// and documents the -file flag.
func TestREADMEDocumentsUsage(t *testing.T) {
	readme := filepath.Join(repoRoot(t), "README.md")
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	content := string(data)
	for _, want := range []string{
		"go run ./cmd/todo add",
		"go run ./cmd/todo list",
		"go run ./cmd/todo done",
		"go run ./cmd/todo rm",
		"-file",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("README.md missing: %s", want)
		}
	}
}
