package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/muli-cohen/pdlc2/semver"
)

func main() {
	sortMode := flag.Bool("sort", false, "read versions from stdin and sort ascending")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	if *sortMode {
		if flag.NArg() > 0 {
			fmt.Fprintln(os.Stderr, "error: -sort reads from stdin and takes no version arguments")
			os.Exit(1)
		}
		runSort()
	} else {
		if flag.NArg() != 2 {
			fmt.Fprintf(os.Stderr, "error: expected 2 version arguments, got %d\n", flag.NArg())
			os.Exit(1)
		}
		runCompare(flag.Arg(0), flag.Arg(1))
	}
}

func runCompare(a, b string) {
	va, err := semver.Parse(a)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	vb, err := semver.Parse(b)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(semver.Compare(va, vb))
}

func runSort() {
	scanner := bufio.NewScanner(os.Stdin)
	var versions []semver.Version
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		v, err := semver.Parse(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		versions = append(versions, v)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error: reading stdin:", err)
		os.Exit(1)
	}
	semver.Sort(versions)
	for _, v := range versions {
		fmt.Println(v.String())
	}
}
