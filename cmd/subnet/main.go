package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/muli-cohen/pdlc2/subnet"
)

func main() {
	jsonFlag := flag.Bool("json", false, "output as JSON")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: subnet <CIDR>")
	}
	flag.Parse()

	switch flag.NArg() {
	case 0:
		fmt.Fprintln(os.Stderr, "Error: CIDR argument required. Usage: subnet <CIDR>")
		os.Exit(1)
	case 1:
		// proceed
	default:
		fmt.Fprintf(os.Stderr, "Error: unexpected extra argument %q\n", flag.Arg(1))
		os.Exit(1)
	}

	n, err := subnet.Parse(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if *jsonFlag {
		data, err := json.Marshal(&n)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		return
	}

	fmt.Printf("Address: %s\n", n.Address)
	fmt.Printf("Network: %s\n", n.NetworkAddress)
	if n.Broadcast != "" {
		fmt.Printf("Broadcast: %s\n", n.Broadcast)
	}
	fmt.Printf("Mask: %s\n", n.Mask)
	if n.FirstHost != "" {
		fmt.Printf("FirstHost: %s\n", n.FirstHost)
	}
	if n.LastHost != "" {
		fmt.Printf("LastHost: %s\n", n.LastHost)
	}
	fmt.Printf("PrefixLength: %d\n", n.PrefixLength)
	fmt.Printf("UsableHosts: %d\n", n.UsableHosts)
}
