package cli

import (
	"coastal-geometry/coastline"
	"coastal-geometry/koch"
	"flag"
	"fmt"
	"io"
	"strings"
)

const (
	defaultOutputDir = "output"
	cmdAll           = "all"
	cmdCoastline     = "coastline"
	cmdParadox       = "paradox"
	cmdKoch          = "koch"
	cmdDimension     = "dimension"
)

type config struct {
	Command    string
	InputPath  string
	OutputPath string
	Iterations int
}

func parseConfig(args []string, stdout, stderr io.Writer) (config, error) {
	if len(args) == 0 {
		printRootUsage(stdout)
		return config{}, flag.ErrHelp
	}

	command := args[0]
	switch command {
	case "-h", "--help", "help":
		printRootUsage(stdout)
		return config{}, flag.ErrHelp
	case cmdAll, cmdCoastline, cmdParadox, cmdKoch, cmdDimension:
	default:
		printRootUsage(stderr)
		return config{}, fmt.Errorf("unknown command %q", command)
	}

	cfg := config{Command: command}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(stderr)

	switch command {
	case cmdAll:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to coastline JSON file")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", koch.MaxIterations, fmt.Sprintf("maximum Koch iterations (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdCoastline:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to coastline JSON file")
		fs.StringVar(&cfg.OutputPath, "output", "", "output SVG path or directory (default: ./output)")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdParadox:
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdKoch:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to coastline JSON file")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", koch.MaxIterations, fmt.Sprintf("maximum Koch iterations (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdDimension:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to coastline JSON file")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", koch.MaxIterations, fmt.Sprintf("maximum Koch iterations (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	}

	if err := fs.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return config{}, err
		}
		return config{}, err
	}

	if fs.NArg() > 0 {
		fs.Usage()
		return config{}, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	if commandUsesIterations(command) && (cfg.Iterations < 0 || cfg.Iterations > koch.MaxIterations) {
		return config{}, fmt.Errorf("iterations must be between 0 and %d", koch.MaxIterations)
	}

	return cfg, nil
}

func commandNeedsCoastline(command string) bool {
	switch command {
	case cmdAll, cmdCoastline, cmdKoch, cmdDimension:
		return true
	default:
		return false
	}
}

func commandUsesIterations(command string) bool {
	switch command {
	case cmdAll, cmdKoch, cmdDimension:
		return true
	default:
		return false
	}
}

func isHelp(err error) bool {
	return err == flag.ErrHelp
}
