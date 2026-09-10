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

## Word Frequency Counter

Count the most frequent words in a file or text stream.

```sh
# Top 10 words in a file
go run ./cmd/wordfreq README.md

# Top 5 words from stdin
cat README.md | go run ./cmd/wordfreq -n 5
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `-n` | 10 | Number of top words to print |

## To-Do List

A persistent command-line to-do list that stores items in a JSON file.

```sh
# Add an item (creates todo.json in the working directory)
go run ./cmd/todo add "buy milk"
# Output: Added item 1

# List all items
go run ./cmd/todo list
# Output: 1 [ ] buy milk

# Mark item 1 as done
go run ./cmd/todo done 1

# List again to see the updated state
go run ./cmd/todo list
# Output: 1 [x] buy milk

# Remove item 1
go run ./cmd/todo rm 1

# List after removal (no output - list is empty)
go run ./cmd/todo list
```

Use the `-file` flag to store items in a custom path:

```sh
go run ./cmd/todo -file /tmp/work.json add "finish report"
go run ./cmd/todo -file /tmp/work.json list
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `-file` | `todo.json` | Path to the JSON storage file |

## Roman Numeral Converter

Convert between integers and Roman numerals. The CLI auto-detects direction: an integer argument converts to a numeral; a numeral argument converts to an integer.

```sh
# Integer to numeral
go run ./cmd/roman 1994
# Output: MCMXCIV

# Numeral to integer
go run ./cmd/roman MCMXCIV
# Output: 1994

# Case-insensitive numeral input
go run ./cmd/roman mcmxciv
# Output: 1994
```

**Error cases:**

```sh
# Out-of-range integer (supported range: 1-3999)
go run ./cmd/roman 0
# stderr: value 0 out of range: roman numerals support 1-3999
# exit code: 1

go run ./cmd/roman 4000
# stderr: value 4000 out of range: roman numerals support 1-3999
# exit code: 1

# Non-canonical numeral
go run ./cmd/roman IIII
# stderr: invalid roman numeral "IIII": not in canonical form
# exit code: 1
```

The library package can also be imported directly:

```go
import "github.com/muli-cohen/pdlc2/roman"

s, err := roman.ToRoman(1994)   // "MCMXCIV", nil
n, err := roman.FromRoman("IV") // 4, nil
```

## Caesar Cipher

Encode or decode text using a Caesar (rotation) cipher.

```sh
# Encode
go run ./cmd/caesar -shift 3 "attack at dawn"
# Output: dwwdfn dw gdzq

# Decode
go run ./cmd/caesar -decode -shift 3 "dwwdfn dw gdzq"
# Output: attack at dawn

# Stdin
echo "abc" | go run ./cmd/caesar
# Output: def
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `-shift` | 3 | Number of positions to shift (any integer; normalised mod 26) |
| `-decode` | false | Decode instead of encode |

Supplying both a positional argument and piped stdin exits non-zero with an error message.
