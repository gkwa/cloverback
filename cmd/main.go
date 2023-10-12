package main

import (
	"flag"
	"log/slog"

	"github.com/taylormonacelli/cloverback"
	"github.com/taylormonacelli/goldbug"
)

var (
	verbose, noExpunge bool
	logFormat          string
)

func main() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output (shorthand)")
	flag.BoolVar(&noExpunge, "no-expunge", false, "Don't expunge data from pushbullet")
	flag.StringVar(&logFormat, "log-format", "", "Log format (text or json)")

	flag.Parse()

	if verbose || logFormat != "" {
		if logFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}
	}

	cloverback.Main(noExpunge)
}
