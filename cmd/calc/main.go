package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/muli-cohen/pdlc2/calc"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: calc <a> <op> <b>  (op: + - * /)")
		os.Exit(1)
	}

	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid operand %q: not a number\n", os.Args[1])
		os.Exit(1)
	}

	b, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid operand %q: not a number\n", os.Args[3])
		os.Exit(1)
	}

	op := os.Args[2]
	var result float64

	switch op {
	case "+":
		result = calc.Add(a, b)
	case "-":
		result = calc.Subtract(a, b)
	case "*":
		result = calc.Multiply(a, b)
	case "/":
		result, err = calc.Divide(a, b)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unsupported operator: %s\n", op)
		os.Exit(1)
	}

	fmt.Printf("%g\n", result)
}
