package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/muli-cohen/pdlc2/tictactoe"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	board := &tictactoe.Board{}
	player := 'X'

	for {
		fmt.Print(board.String())
		fmt.Printf("Player %c, enter move (row col): ", player)

		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "Error: unexpected end of input.")
			os.Exit(1)
		}
		line := scanner.Text()

		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Println("Invalid input. Please enter row and column as two integers (0-2).")
			continue
		}

		row, err1 := strconv.Atoi(parts[0])
		col, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid input. Please enter row and column as two integers (0-2).")
			continue
		}

		if row < 0 || row > 2 || col < 0 || col > 2 {
			fmt.Println("Invalid move: row and column must be between 0 and 2.")
			continue
		}

		if err := board.Move(row, col, player); err != nil {
			fmt.Println("Invalid move: cell already occupied.")
			continue
		}

		if w := board.Winner(); w != 0 {
			fmt.Print(board.String())
			fmt.Printf("Player %c wins!\n", w)
			os.Exit(0)
		}

		if board.IsDraw() {
			fmt.Print(board.String())
			fmt.Println("It's a draw!")
			os.Exit(0)
		}

		if player == 'X' {
			player = 'O'
		} else {
			player = 'X'
		}
	}
}
