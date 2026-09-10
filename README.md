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
