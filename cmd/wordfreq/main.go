package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/muli-cohen/pdlc2/wordfreq"
)

func main() {
	n := flag.Int("n", 10, "number of top words to print")
	flag.Parse()

	if *n < 1 {
		fmt.Fprintln(os.Stderr, "error: -n must be at least 1")
		os.Exit(1)
	}
	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "usage: wordfreq [-n N] [file]")
		os.Exit(1)
	}

	r := os.Stdin
	if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot open file %q: %v\n", flag.Arg(0), err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	counts, err := wordfreq.Count(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	for _, e := range wordfreq.Top(counts, *n) {
		fmt.Printf("%d %s\n", e.Count, e.Word)
	}
}
