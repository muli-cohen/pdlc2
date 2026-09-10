package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "tictactoe-e2e-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create temp dir:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	binary = filepath.Join(dir, "tictactoe")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n%s\n", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func gameRun(input string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func TestXWinsRow(t *testing.T) {
	// X(0,0) O(1,0) X(0,1) O(1,1) X(0,2) -> X fills row 0
	stdout, _, code := gameRun("0 0\n1 0\n0 1\n1 1\n0 2\n")
	if code != 0 {
		t.Errorf("exit code = %d; want 0", code)
	}
	if !strings.Contains(stdout, "Player X wins!") {
		t.Errorf("stdout = %q; want to contain 'Player X wins!'", stdout)
	}
}

func TestOWinsCol(t *testing.T) {
	// O(0,2) X(0,0) O(1,2) X(1,0) O(2,2) -> O fills col 2
	stdout, _, code := gameRun("0 0\n0 2\n1 0\n1 2\n0 1\n2 2\n")
	if code != 0 {
		t.Errorf("exit code = %d; want 0", code)
	}
	if !strings.Contains(stdout, "Player O wins!") {
		t.Errorf("stdout = %q; want to contain 'Player O wins!'", stdout)
	}
}

func TestDraw(t *testing.T) {
	// Board: X O X / X X O / O X O  — full, no winner
	// Move sequence (alternating X/O):
	// X(0,0) O(0,1) X(0,2) O(1,2) X(1,0) O(2,0) X(1,1) O(2,2) X(2,1)
	stdout, _, code := gameRun("0 0\n0 1\n0 2\n1 2\n1 0\n2 0\n1 1\n2 2\n2 1\n")
	if code != 0 {
		t.Errorf("exit code = %d; want 0", code)
	}
	if !strings.Contains(stdout, "It's a draw!") {
		t.Errorf("stdout = %q; want to contain \"It's a draw!\"", stdout)
	}
}

func TestEOFExitsNonZero(t *testing.T) {
	_, stderr, code := gameRun("")
	if code == 0 {
		t.Error("exit code = 0 on EOF; want non-zero")
	}
	if !strings.Contains(stderr, "Error: unexpected end of input.") {
		t.Errorf("stderr = %q; want 'Error: unexpected end of input.'", stderr)
	}
}

func TestInvalidInputRePrompts(t *testing.T) {
	// non-integer tokens followed by EOF
	stdout, _, _ := gameRun("abc\n")
	if !strings.Contains(stdout, "Invalid input.") {
		t.Errorf("stdout = %q; want to contain 'Invalid input.'", stdout)
	}
}

func TestOutOfRangeRePrompts(t *testing.T) {
	// coordinate outside [0,2] followed by EOF
	stdout, _, _ := gameRun("5 0\n")
	if !strings.Contains(stdout, "Invalid move: row and column must be between 0 and 2.") {
		t.Errorf("stdout = %q; want out-of-range error message", stdout)
	}
}

func TestOccupiedCellRePrompts(t *testing.T) {
	// move on occupied cell followed by EOF
	stdout, _, _ := gameRun("0 0\n0 0\n")
	if !strings.Contains(stdout, "Invalid move: cell already occupied.") {
		t.Errorf("stdout = %q; want occupied-cell error message", stdout)
	}
}

func TestInitialBoardAndPrompt(t *testing.T) {
	// Immediate EOF: should still see the empty board and first prompt
	stdout, _, _ := gameRun("")
	if !strings.Contains(stdout, "Player X, enter move (row col): ") {
		t.Errorf("stdout = %q; want initial prompt for player X", stdout)
	}
}
