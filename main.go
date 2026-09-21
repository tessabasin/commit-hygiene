package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: commit-hygiene [logfile]")
		os.Exit(2)
	}

	input := os.Stdin
	if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit-hygiene: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	commits, err := parseLog(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "commit-hygiene: %v\n", err)
		os.Exit(1)
	}
	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "commit-hygiene: no commits found in input")
		os.Exit(1)
	}

	buildReport(commits).print(os.Stdout)
}
