package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options] <file-a> <file-b>\n\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	files, err := loader.LoadMultiple(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	results := diff.Compare(files[0].Vars, files[1].Vars)

	if err := reporter.Report(os.Stdout, results, *format); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if len(results) > 0 {
		os.Exit(2)
	}
}
