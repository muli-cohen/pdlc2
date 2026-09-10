package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/muli-cohen/pdlc2/password"
)

func main() {
	length := flag.Int("length", 16, "password length (8-128)")
	upper := flag.Bool("upper", true, "include uppercase letters")
	digits := flag.Bool("digits", true, "include digits")
	symbols := flag.Bool("symbols", true, "include symbols")
	flag.Parse()

	opts := password.Options{
		Lowercase: true,
		Uppercase: *upper,
		Digits:    *digits,
		Symbols:   *symbols,
	}
	pw, err := password.Generate(*length, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(pw)
}
