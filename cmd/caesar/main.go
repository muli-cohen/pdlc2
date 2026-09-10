package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/muli-cohen/pdlc2/caesar"
)

func main() {
	shift := flag.Int("shift", 3, "shift amount")
	decode := flag.Bool("decode", false, "decode instead of encode")
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
		fmt.Fprintln(os.Stderr, "Error: cannot accept both a positional argument and piped stdin")
		os.Exit(1)
	case flag.NArg() > 0:
		input = flag.Arg(0)
	case stdinPiped:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: reading stdin:", err)
			os.Exit(1)
		}
		input = strings.TrimSuffix(string(data), "\n")
	default:
		fmt.Fprintln(os.Stderr, "usage: caesar [-shift N] [-decode] <text>")
		fmt.Fprintln(os.Stderr, "       echo <text> | caesar [-shift N] [-decode]")
		os.Exit(1)
	}

	var result string
	if *decode {
		result, _ = caesar.Decode(input, *shift)
	} else {
		result, _ = caesar.Encode(input, *shift)
	}
	fmt.Println(result)
}
