package tictactoe

import "testing"

func TestWinRow0(t *testing.T) {
	var b Board
	b.Move(0, 0, 'X')
	b.Move(0, 1, 'X')
	b.Move(0, 2, 'X')
	if w := b.Winner(); w != 'X' {
		t.Errorf("Winner() = %c; want X", w)
	}
}

func TestWinRow1(t *testing.T) {
	var b Board
	b.Move(1, 0, 'O')
	b.Move(1, 1, 'O')
	b.Move(1, 2, 'O')
	if w := b.Winner(); w != 'O' {
		t.Errorf("Winner() = %c; want O", w)
	}
}

func TestWinRow2(t *testing.T) {
	var b Board
	b.Move(2, 0, 'X')
	b.Move(2, 1, 'X')
	b.Move(2, 2, 'X')
	if w := b.Winner(); w != 'X' {
		t.Errorf("Winner() = %c; want X", w)
	}
}

func TestWinCol0(t *testing.T) {
	var b Board
	b.Move(0, 0, 'O')
	b.Move(1, 0, 'O')
	b.Move(2, 0, 'O')
	if w := b.Winner(); w != 'O' {
		t.Errorf("Winner() = %c; want O", w)
	}
}

func TestWinCol1(t *testing.T) {
	var b Board
	b.Move(0, 1, 'X')
	b.Move(1, 1, 'X')
	b.Move(2, 1, 'X')
	if w := b.Winner(); w != 'X' {
		t.Errorf("Winner() = %c; want X", w)
	}
}

func TestWinCol2(t *testing.T) {
	var b Board
	b.Move(0, 2, 'O')
	b.Move(1, 2, 'O')
	b.Move(2, 2, 'O')
	if w := b.Winner(); w != 'O' {
		t.Errorf("Winner() = %c; want O", w)
	}
}

func TestWinDiagMain(t *testing.T) {
	var b Board
	b.Move(0, 0, 'X')
	b.Move(1, 1, 'X')
	b.Move(2, 2, 'X')
	if w := b.Winner(); w != 'X' {
		t.Errorf("Winner() = %c; want X", w)
	}
}

func TestWinDiagAnti(t *testing.T) {
	var b Board
	b.Move(0, 2, 'O')
	b.Move(1, 1, 'O')
	b.Move(2, 0, 'O')
	if w := b.Winner(); w != 'O' {
		t.Errorf("Winner() = %c; want O", w)
	}
}

func TestDraw(t *testing.T) {
	var b Board
	// X O X
	// X X O
	// O X O
	moves := [][3]int{{0, 0, 'X'}, {0, 1, 'O'}, {0, 2, 'X'}, {1, 0, 'X'}, {1, 1, 'X'}, {1, 2, 'O'}, {2, 0, 'O'}, {2, 1, 'X'}, {2, 2, 'O'}}
	for _, m := range moves {
		if err := b.Move(m[0], m[1], rune(m[2])); err != nil {
			t.Fatalf("Move(%d,%d,%c) unexpected error: %v", m[0], m[1], m[2], err)
		}
	}
	if w := b.Winner(); w != 0 {
		t.Errorf("Winner() = %c; want 0 (no winner)", w)
	}
	if !b.IsDraw() {
		t.Error("IsDraw() = false; want true")
	}
}

func TestMoveOccupied(t *testing.T) {
	var b Board
	b.Move(1, 1, 'X')
	if err := b.Move(1, 1, 'O'); err == nil {
		t.Error("Move to occupied cell returned nil error; want non-nil")
	}
}

func TestMoveRowOutOfRange(t *testing.T) {
	var b Board
	if err := b.Move(-1, 0, 'X'); err == nil {
		t.Error("Move(-1,0,'X') returned nil error; want non-nil")
	}
	if err := b.Move(3, 0, 'X'); err == nil {
		t.Error("Move(3,0,'X') returned nil error; want non-nil")
	}
}

func TestMoveColOutOfRange(t *testing.T) {
	var b Board
	if err := b.Move(0, -1, 'X'); err == nil {
		t.Error("Move(0,-1,'X') returned nil error; want non-nil")
	}
	if err := b.Move(0, 3, 'X'); err == nil {
		t.Error("Move(0,3,'X') returned nil error; want non-nil")
	}
}
