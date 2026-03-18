package cli

import (
	"coastal-geometry/internal/domain/coastline"
	"coastal-geometry/internal/domain/generators/koch"
	"flag"
	"fmt"
	"io"
	"strings"
)

const (
	defaultOutputDir = "output"
	cmdReal          = "real"
	cmdModel         = "model"
	cmdSource        = "source"
	cmdAll           = "all"
	cmdCoastline     = "coastline"
	cmdParadox       = "paradox"
	cmdKoch          = "koch"
	cmdKochOrganic   = "koch-organic"
	cmdDimension     = "dimension"
)

type config struct {
	Command      string
	InputPath    string
	SourceURL    string
	Refresh      bool
	OutputPath   string
	Iterations   int
	Seed         int64
	AngleJitter  float64
	HeightJitter float64
}

func parseConfig(args []string, stdout, stderr io.Writer) (config, error) {
	if len(args) == 0 {
		printRootUsage(stdout)
		return config{}, flag.ErrHelp
	}

	if isHelpToken(args[0]) {
		printRootUsage(stdout)
		return config{}, flag.ErrHelp
	}

	command, commandArgs, err := resolveCommand(args, stdout, stderr)
	if err != nil {
		return config{}, err
	}

	cfg := config{Command: command}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(stderr)

	switch command {
	case cmdSource:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before saving a snapshot")
		fs.StringVar(&cfg.OutputPath, "output", "", "snapshot file or directory (default: ./data/snapshots)")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdAll:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before running")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum organic Koch iterations (0-%d)", koch.MaxIterations))
		fs.Int64Var(&cfg.Seed, "seed", 42, "random seed for organic coastline generation")
		fs.Float64Var(&cfg.AngleJitter, "angle-jitter", 18, "maximum random angle deviation in degrees")
		fs.Float64Var(&cfg.HeightJitter, "height-jitter", 0.25, "maximum random height deviation as a ratio")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdCoastline:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before running")
		fs.StringVar(&cfg.OutputPath, "output", "", "output SVG path or directory (default: ./output)")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdParadox:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before running")
		fs.IntVar(&cfg.Iterations, "iterations", 4, fmt.Sprintf("maximum paradox detail levels (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdKoch:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before running")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum Koch iterations (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdKochOrganic:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before running")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum organic Koch iterations (0-%d)", koch.MaxIterations))
		fs.Int64Var(&cfg.Seed, "seed", 42, "random seed for organic coastline generation")
		fs.Float64Var(&cfg.AngleJitter, "angle-jitter", 18, "maximum random angle deviation in degrees")
		fs.Float64Var(&cfg.HeightJitter, "height-jitter", 0.25, "maximum random height deviation as a ratio")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdDimension:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.BoolVar(&cfg.Refresh, "refresh", false, "force refresh of the remote GeoJSON cache before running")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum organic Koch iterations (0-%d)", koch.MaxIterations))
		fs.Int64Var(&cfg.Seed, "seed", 42, "random seed for organic coastline generation")
		fs.Float64Var(&cfg.AngleJitter, "angle-jitter", 18, "maximum random angle deviation in degrees")
		fs.Float64Var(&cfg.HeightJitter, "height-jitter", 0.25, "maximum random height deviation as a ratio")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	}

	if err := fs.Parse(commandArgs); err != nil {
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
	if command == cmdAll || command == cmdKochOrganic || command == cmdDimension {
		if cfg.AngleJitter < 0 {
			return config{}, fmt.Errorf("angle-jitter must be non-negative")
		}
		if cfg.HeightJitter < 0 {
			return config{}, fmt.Errorf("height-jitter must be non-negative")
		}
	}

	return cfg, nil
}

func commandNeedsCoastline(command string) bool {
	switch command {
	case cmdAll, cmdCoastline, cmdParadox, cmdKoch, cmdKochOrganic, cmdDimension:
		return true
	default:
		return false
	}
}

func commandUsesIterations(command string) bool {
	switch command {
	case cmdAll, cmdParadox, cmdKoch, cmdKochOrganic, cmdDimension:
		return true
	default:
		return false
	}
}

func isHelp(err error) bool {
	return err == flag.ErrHelp
}

func resolveCommand(args []string, stdout, stderr io.Writer) (string, []string, error) {
	switch args[0] {
	case cmdReal:
		return resolveGroupedCommand(cmdReal, args[1:], stdout, stderr)
	case cmdModel:
		return resolveGroupedCommand(cmdModel, args[1:], stdout, stderr)
	case cmdSource, cmdAll, cmdCoastline, cmdParadox, cmdKoch, cmdKochOrganic, cmdDimension:
		return args[0], args[1:], nil
	default:
		printRootUsage(stderr)
		return "", nil, fmt.Errorf("unknown command %q", args[0])
	}
}

func resolveGroupedCommand(group string, args []string, stdout, stderr io.Writer) (string, []string, error) {
	if len(args) == 0 || isHelpToken(args[0]) {
		printGroupUsage(stdout, group)
		return "", nil, flag.ErrHelp
	}

	command := args[0]
	if !commandBelongsToGroup(command, group) {
		printGroupUsage(stderr, group)
		return "", nil, fmt.Errorf("unknown %s command %q", group, command)
	}

	return command, args[1:], nil
}

func commandBelongsToGroup(command, group string) bool {
	switch group {
	case cmdReal:
		return command == cmdCoastline
	case cmdModel:
		switch command {
		case cmdParadox, cmdKoch, cmdKochOrganic, cmdDimension:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func isHelpToken(arg string) bool {
	switch arg {
	case "-h", "--help", "help":
		return true
	default:
		return false
	}
}
