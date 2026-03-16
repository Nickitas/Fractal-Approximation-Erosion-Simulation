package cli

import (
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

	if err := executeCommand(app); err != nil {
		exitWithError(stderr, err)
	}
}

func exitWithError(stderr io.Writer, err error) {
	fmt.Fprintf(stderr, "error: %v\n", err)
	os.Exit(1)
}
