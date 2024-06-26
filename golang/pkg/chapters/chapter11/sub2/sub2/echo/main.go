package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	n = flag.Bool("n", false, "omit trailing newline")
	s = flag.String("s", " ", "separator")

	// Modified during testing
	out io.Writer = os.Stdout
)

func main() {
	flag.Parse()
	if err := echo(!*n, *s, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "echo %v\n", err)
	}
}

func echo(newline bool, sep string, args []string) error {
	fmt.Fprint(out, strings.Join(args, sep))
	if newline {
		fmt.Fprintln(out)
	}
	return nil
}
