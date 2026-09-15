package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/muli-cohen/pdlc2/rle"
)

func main() {
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
		fmt.Fprintln(os.Stderr, "Error: provide input as an argument or via stdin, not both")
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
		fmt.Fprintln(os.Stderr, "usage: rle [-decode] <text>")
		fmt.Fprintln(os.Stderr, "       echo <text> | rle [-decode]")
		os.Exit(1)
	}

	var result string
	if *decode {
		result, err = rle.Decode(input)
	} else {
		result, err = rle.Encode(input)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
