package tictactoe

import (
	"fmt"
	"strings"
)

type Board struct {
	cells [3][3]rune
}

func (b *Board) Move(row, col int, player rune) error {
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return fmt.Errorf("row and column must be between 0 and 2")
	}
	if b.cells[row][col] != 0 {
		return fmt.Errorf("cell already occupied")
	}
	b.cells[row][col] = player
	return nil
}

func (b *Board) Winner() rune {
	for r := 0; r < 3; r++ {
		if b.cells[r][0] != 0 && b.cells[r][0] == b.cells[r][1] && b.cells[r][1] == b.cells[r][2] {
			return b.cells[r][0]
		}
	}
	for c := 0; c < 3; c++ {
		if b.cells[0][c] != 0 && b.cells[0][c] == b.cells[1][c] && b.cells[1][c] == b.cells[2][c] {
			return b.cells[0][c]
		}
	}
	if b.cells[0][0] != 0 && b.cells[0][0] == b.cells[1][1] && b.cells[1][1] == b.cells[2][2] {
		return b.cells[0][0]
	}
	if b.cells[0][2] != 0 && b.cells[0][2] == b.cells[1][1] && b.cells[1][1] == b.cells[2][0] {
		return b.cells[0][2]
	}
	return 0
}

func (b *Board) IsDraw() bool {
	if b.Winner() != 0 {
		return false
	}
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if b.cells[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

func displayRune(r rune) rune {
	if r == 0 {
		return ' '
	}
	return r
}

func (b *Board) String() string {
	var sb strings.Builder
	for r := 0; r < 3; r++ {
		fmt.Fprintf(&sb, " %c | %c | %c\n",
			displayRune(b.cells[r][0]),
			displayRune(b.cells[r][1]),
			displayRune(b.cells[r][2]),
		)
		if r < 2 {
			sb.WriteString("-----------\n")
		}
	}
	return sb.String()
}
