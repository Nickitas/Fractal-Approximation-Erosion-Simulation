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

	command := args[0]
	switch command {
	case "-h", "--help", "help":
		printRootUsage(stdout)
		return config{}, flag.ErrHelp
	case cmdAll, cmdCoastline, cmdParadox, cmdKoch, cmdKochOrganic, cmdDimension:
	default:
		printRootUsage(stderr)
		return config{}, fmt.Errorf("unknown command %q", command)
	}

	cfg := config{Command: command}
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(stderr)

	switch command {
	case cmdAll:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum organic Koch iterations (0-%d)", koch.MaxIterations))
		fs.Int64Var(&cfg.Seed, "seed", 42, "random seed for organic coastline generation")
		fs.Float64Var(&cfg.AngleJitter, "angle-jitter", 18, "maximum random angle deviation in degrees")
		fs.Float64Var(&cfg.HeightJitter, "height-jitter", 0.25, "maximum random height deviation as a ratio")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdCoastline:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.StringVar(&cfg.OutputPath, "output", "", "output SVG path or directory (default: ./output)")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdParadox:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.IntVar(&cfg.Iterations, "iterations", 4, fmt.Sprintf("maximum paradox detail levels (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdKoch:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum Koch iterations (0-%d)", koch.MaxIterations))
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdKochOrganic:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum organic Koch iterations (0-%d)", koch.MaxIterations))
		fs.Int64Var(&cfg.Seed, "seed", 42, "random seed for organic coastline generation")
		fs.Float64Var(&cfg.AngleJitter, "angle-jitter", 18, "maximum random angle deviation in degrees")
		fs.Float64Var(&cfg.HeightJitter, "height-jitter", 0.25, "maximum random height deviation as a ratio")
		fs.Usage = func() { printCommandUsage(stdout, command) }
	case cmdDimension:
		fs.StringVar(&cfg.InputPath, "input", coastline.DefaultCoastlineJSONPath, "path to local coastline JSON/GeoJSON fallback file")
		fs.StringVar(&cfg.SourceURL, "source-url", coastline.DefaultCoastlineGeoJSONURL, "remote GeoJSON URL for coastline data; empty string disables HTTP loading")
		fs.StringVar(&cfg.OutputPath, "output", "", "output directory for generated visualizations (default: ./output)")
		fs.IntVar(&cfg.Iterations, "iterations", 5, fmt.Sprintf("maximum organic Koch iterations (0-%d)", koch.MaxIterations))
		fs.Int64Var(&cfg.Seed, "seed", 42, "random seed for organic coastline generation")
		fs.Float64Var(&cfg.AngleJitter, "angle-jitter", 18, "maximum random angle deviation in degrees")
		fs.Float64Var(&cfg.HeightJitter, "height-jitter", 0.25, "maximum random height deviation as a ratio")
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
