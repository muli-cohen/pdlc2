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

## Joke Teller

Print a random joke:

```sh
go run ./cmd/joke
```

Running the command repeatedly will produce different jokes. No arguments are accepted; passing any argument exits with code 1 and prints a usage message to stderr.

## Password Generator

Generate cryptographically secure passwords:

```sh
# Default 16-character password (uppercase, digits, and symbols included)
go run ./cmd/passgen

# 20-character password
go run ./cmd/passgen -length 20

# Password without symbols
go run ./cmd/passgen -symbols=false

# Error: length below minimum
go run ./cmd/passgen -length 4
# stderr: error: length must be between 8 and 128
# exit code: 1
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `-length` | 16 | Password length (8-128) |
| `-upper` | true | Include uppercase letters |
| `-digits` | true | Include digits |
| `-symbols` | true | Include symbols |

Lowercase letters are always included. Two consecutive runs with the same flags will produce different passwords.
