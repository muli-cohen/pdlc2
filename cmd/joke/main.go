package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/muli-cohen/pdlc2/jokes"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintln(os.Stderr, "Usage: joke\nNo arguments are accepted.")
		os.Exit(1)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(jokes.Random(r))
}
