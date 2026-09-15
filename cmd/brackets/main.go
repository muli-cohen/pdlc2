package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/muli-cohen/pdlc2/brackets"
)

func main() {
	position := flag.Bool("position", false, "print the 0-based byte index of the first mismatch")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot stat stdin:", err)
		os.Exit(1)
	}
	stdinPiped := (stat.Mode() & os.ModeCharDevice) == 0

	var input string
	switch {
	case flag.NArg() > 0 && stdinPiped:
		fmt.Fprintln(os.Stderr, "error: supply input as an argument or via stdin, not both")
		os.Exit(1)
	case flag.NArg() > 0:
		input = flag.Arg(0)
	case stdinPiped:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: reading stdin:", err)
			os.Exit(1)
		}
		input = strings.TrimRight(string(data), "\r\n")
	default:
		fmt.Fprintln(os.Stderr, "usage: brackets [-position] <text>")
		fmt.Fprintln(os.Stderr, "       echo <text> | brackets [-position]")
		os.Exit(1)
	}

	if *position {
		idx, err := brackets.FirstMismatch(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(idx)
		return
	}

	if brackets.IsBalanced(input) {
		fmt.Println("balanced")
	} else {
		fmt.Println("unbalanced")
	}
}
