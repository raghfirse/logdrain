package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/logdrain/internal/filter"
	"github.com/user/logdrain/internal/formatter"
	"github.com/user/logdrain/internal/source"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		rawFormat  = flag.Bool("raw", false, "output raw JSON lines without formatting")
		filterExpr = flag.String("filter", "", "comma-separated filter expressions (e.g. level=error,service=api)")
	)
	flag.Parse()

	files := flag.Args()

	var sources []*source.Source
	if len(files) == 0 {
		src, err := source.NewFromStdin()
		if err != nil {
			return fmt.Errorf("stdin: %w", err)
		}
		sources = append(sources, src)
	} else {
		for _, f := range files {
			src, err := source.NewFromFile(f)
			if err != nil {
				return fmt.Errorf("open %s: %w", f, err)
			}
			sources = append(sources, src)
		}
	}

	var rules []filter.Rule
	if *filterExpr != "" {
		for _, expr := range strings.Split(*filterExpr, ",") {
			expr = strings.TrimSpace(expr)
			if expr == "" {
				continue
			}
			f, err := filter.New(expr)
			if err != nil {
				return fmt.Errorf("invalid filter %q: %w", expr, err)
			}
			rules = append(rules, f...)
		}
	}

	fmt_mode := formatter.Pretty
	if *rawFormat {
		fmt_mode = formatter.Raw
	}
	fmt := formatter.New(os.Stdout, fmt_mode)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	merged := source.Merge(ctx, sources...)
	for line := range merged {
		if len(rules) > 0 && !filter.Match(rules, line.Text) {
			continue
		}
		if err := fmt.Write(line.Source, line.Text); err != nil {
			return err
		}
	}
	return nil
}
