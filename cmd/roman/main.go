package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/muli-cohen/pdlc2/roman"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: roman <integer|numeral>")
		os.Exit(1)
	}

	arg := os.Args[1]
	if n, err := strconv.Atoi(arg); err == nil {
		s, err := roman.ToRoman(n)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(s)
	} else {
		n, err := roman.FromRoman(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(n)
	}
}
