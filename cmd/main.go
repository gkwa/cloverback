package main

import (
	"flag"
	"log/slog"

	"github.com/taylormonacelli/cloverback"
)

var (
	verbose   bool
	logFormat string
)

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
	flag.StringVar(&logFormat, "log-format", "", "Log format (text or json)")

	flag.Parse()

	if verbose || logFormat != "" {
		if logFormat == "json" {
			cloverback.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			cloverback.SetDefaultLoggerText(slog.LevelDebug)
		}
	}

	cloverback.Main()
}
