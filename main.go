package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	maxSubjectLen := flag.Int("max-subject-len", 50, "flag subject lines longer than this many characters")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: commit-hygiene [-max-subject-len N] [logfile]")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		flag.Usage()
		os.Exit(2)
	}
	if *maxSubjectLen < 1 {
		fmt.Fprintln(os.Stderr, "commit-hygiene: -max-subject-len must be at least 1")
		os.Exit(2)
	}

	input := os.Stdin
	if len(args) == 1 {
		f, err := os.Open(args[0])
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

	buildReport(commits, *maxSubjectLen).print(os.Stdout)
}
