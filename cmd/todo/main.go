package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/muli-cohen/pdlc2/todo"
)

const usageMsg = `usage: todo [-file <path>] <subcommand> [args]

Subcommands:
  add <title>   add a new item
  done <id>     mark an item as done
  rm <id>       remove an item
  list          list all items

Flags:
  -file string  path to the JSON storage file (default "todo.json")
`

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

func usage() {
	fmt.Fprint(os.Stderr, usageMsg)
	os.Exit(1)
}

func main() {
	file := flag.String("file", "todo.json", "path to the JSON storage file")
	flag.Usage = func() { fmt.Fprint(os.Stderr, usageMsg) }
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
	}

	sub := args[0]
	rest := args[1:]

	switch sub {
	case "add":
		cmdAdd(*file, rest)
	case "done":
		cmdDone(*file, rest)
	case "rm":
		cmdRm(*file, rest)
	case "list":
		cmdList(*file)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown subcommand %q\n", sub)
		fmt.Fprint(os.Stderr, usageMsg)
		os.Exit(1)
	}
}

func cmdAdd(file string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: todo add <title>")
		os.Exit(1)
	}
	l, err := todo.Load(file)
	if err != nil {
		fatal("%v", err)
	}
	item, err := l.Add(args[0])
	if err != nil {
		fatal("%v", err)
	}
	if err := l.Save(file); err != nil {
		fatal("%v", err)
	}
	fmt.Printf("Added item %d\n", item.Id)
}

func cmdDone(file string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: todo done <id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage: todo done <id>\nerror: %q is not a valid id\n", args[0])
		os.Exit(1)
	}
	l, err := todo.Load(file)
	if err != nil {
		fatal("%v", err)
	}
	if err := l.Done(id); err != nil {
		fatal("%v", err)
	}
	if err := l.Save(file); err != nil {
		fatal("%v", err)
	}
}

func cmdRm(file string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: todo rm <id>")
		os.Exit(1)
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "usage: todo rm <id>\nerror: %q is not a valid id\n", args[0])
		os.Exit(1)
	}
	l, err := todo.Load(file)
	if err != nil {
		fatal("%v", err)
	}
	if err := l.Remove(id); err != nil {
		fatal("%v", err)
	}
	if err := l.Save(file); err != nil {
		fatal("%v", err)
	}
}

func cmdList(file string) {
	l, err := todo.Load(file)
	if err != nil {
		fatal("%v", err)
	}
	for _, item := range l.Items() {
		mark := " "
		if item.Done {
			mark = "x"
		}
		fmt.Printf("%d [%s] %s\n", item.Id, mark, item.Title)
	}
}
