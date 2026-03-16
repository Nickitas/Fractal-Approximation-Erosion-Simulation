package cli

import (
	"coastal-geometry/internal/domain/coastline"
	"fmt"
	"io"
	"os"
)

func Run(args []string, stdout, stderr io.Writer) {
	cfg, err := parseConfig(args, stdout, stderr)
	if err != nil {
		if isHelp(err) {
			return
		}
		exitWithError(stderr, err)
	}

	app, err := NewApp(cfg)
	if err != nil {
		exitWithError(stderr, err)
	}

	printValidationReport(stdout, app.Validation)

	if err := executeCommand(app); err != nil {
		exitWithError(stderr, err)
	}
}

func exitWithError(stderr io.Writer, err error) {
	fmt.Fprintf(stderr, "error: %v\n", err)
	os.Exit(1)
}

func printValidationReport(w io.Writer, report coastline.ValidationReport) {
	for _, fix := range report.Fixes {
		fmt.Fprintf(w, "fix: %s\n", fix)
	}
	for _, warning := range report.Warnings {
		fmt.Fprintf(w, "warning: %s\n", warning)
	}
}
