# pdlc2

Seed repository for PDLC autonomous-run testing.

## Usage

```sh
# Addition
go run ./cmd/calc 1.5 + 2.5
# Output: 4

# Subtraction
go run ./cmd/calc 10 - 3
# Output: 7

# Division (whole-number result)
go run ./cmd/calc 6 / 3
# Output: 2

# Fractional result
go run ./cmd/calc 1 / 2
# Output: 0.5

# Multiplication (quote * in most shells)
go run ./cmd/calc 10 '*' 3
# Output: 30
```

**Error cases:**

```sh
# Division by zero
go run ./cmd/calc 1 / 0
# stderr: division by zero
# exit code: 1

# Unsupported operator
go run ./cmd/calc 1 ^ 2
# stderr: unsupported operator: ^
# exit code: 1
```

## Tic Tac Toe

A two-player hot-seat tic tac toe game playable entirely in the terminal.

```sh
go run ./cmd/tictactoe
```

On launch, the empty 3x3 board is printed and player X is prompted first:

```
   |   |  
-----------
   |   |  
-----------
   |   |  
Player X, enter move (row col): 
```

Enter a move as two space-separated integers (row and column, each 0-2). Players alternate X then O. The game announces the winner or a draw and exits with code 0.

**Error cases:**

```sh
# Non-integer input
Player X, enter move (row col): abc
Invalid input. Please enter row and column as two integers (0-2).

# Out-of-range coordinate
Player X, enter move (row col): 3 0
Invalid move: row and column must be between 0 and 2.

# Occupied cell
Player X, enter move (row col): 0 0
Invalid move: cell already occupied.

# Closing stdin mid-game (Ctrl-D / EOF)
# stderr: Error: unexpected end of input.
# exit code: 1
```
