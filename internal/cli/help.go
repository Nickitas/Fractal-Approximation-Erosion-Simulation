package cli

import (
	"coastal-geometry/coastline"
	"coastal-geometry/koch"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func printRootUsage(w io.Writer) {
	bin := filepath.Base(os.Args[0])
	fmt.Fprintf(w, "Usage: %s <command> [flags]\n\n", bin)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintf(w, "  %-10s run every demonstration\n", cmdAll)
	fmt.Fprintf(w, "  %-10s run coastline metrics and save coastline.svg\n", cmdCoastline)
	fmt.Fprintf(w, "  %-10s run coastline paradox demo only\n", cmdParadox)
	fmt.Fprintf(w, "  %-10s run Koch curve demo and save koch_iter_0..N.svg\n", cmdKoch)
	fmt.Fprintf(w, "  %-10s run fractal dimension demo and save koch_iter_0..N.svg\n", cmdDimension)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintf(w, "  %s coastline\n", bin)
	fmt.Fprintf(w, "  %s koch --iterations 4 --output ./output/koch\n", bin)
	fmt.Fprintf(w, "  %s dimension --iterations 6 --input data/black-sea.json\n", bin)
	fmt.Fprintf(w, "  %s all --output ./output/full-run\n", bin)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "Run '%s <command> --help' for command-specific flags.\n", bin)
}

func printCommandUsage(w io.Writer, command string) {
	bin := filepath.Base(os.Args[0])
	switch command {
	case cmdAll:
		fmt.Fprintf(w, "Usage: %s %s [flags]\n\n", bin, cmdAll)
		fmt.Fprintln(w, "Runs every current demonstration: coastline, paradox, Koch, and fractal dimension.")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Flags:")
		fmt.Fprintln(w, "  --input string")
		fmt.Fprintf(w, "        path to coastline JSON file (default %q)\n", coastline.DefaultCoastlineJSONPath)
		fmt.Fprintln(w, "  --iterations int")
		fmt.Fprintf(w, "        maximum Koch iterations (0-%d)\n", koch.MaxIterations)
		fmt.Fprintln(w, "  --output string")
		fmt.Fprintln(w, "        output directory for generated visualizations (default: ./output)")
	case cmdCoastline:
		fmt.Fprintf(w, "Usage: %s %s [flags]\n\n", bin, cmdCoastline)
		fmt.Fprintln(w, "Prints coastline metrics and saves a coastline SVG.")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Flags:")
		fmt.Fprintln(w, "  --input string")
		fmt.Fprintf(w, "        path to coastline JSON file (default %q)\n", coastline.DefaultCoastlineJSONPath)
		fmt.Fprintln(w, "  --output string")
		fmt.Fprintln(w, "        output SVG path or directory (default: ./output)")
	case cmdParadox:
		fmt.Fprintf(w, "Usage: %s %s\n\n", bin, cmdParadox)
		fmt.Fprintln(w, "Runs the coastline paradox demonstration.")
	case cmdKoch:
		fmt.Fprintf(w, "Usage: %s %s [flags]\n\n", bin, cmdKoch)
		fmt.Fprintln(w, "Builds the Koch approximation and saves koch_iter_0..N.svg.")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Flags:")
		fmt.Fprintln(w, "  --input string")
		fmt.Fprintf(w, "        path to coastline JSON file (default %q)\n", coastline.DefaultCoastlineJSONPath)
		fmt.Fprintln(w, "  --iterations int")
		fmt.Fprintf(w, "        maximum Koch iterations (0-%d)\n", koch.MaxIterations)
		fmt.Fprintln(w, "  --output string")
		fmt.Fprintln(w, "        output directory for generated visualizations (default: ./output)")
	case cmdDimension:
		fmt.Fprintf(w, "Usage: %s %s [flags]\n\n", bin, cmdDimension)
		fmt.Fprintln(w, "Computes fractal dimension across Koch iterations and saves koch_iter_0..N.svg.")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Flags:")
		fmt.Fprintln(w, "  --input string")
		fmt.Fprintf(w, "        path to coastline JSON file (default %q)\n", coastline.DefaultCoastlineJSONPath)
		fmt.Fprintln(w, "  --iterations int")
		fmt.Fprintf(w, "        maximum Koch iterations (0-%d)\n", koch.MaxIterations)
		fmt.Fprintln(w, "  --output string")
		fmt.Fprintln(w, "        output directory for generated visualizations (default: ./output)")
	}
}
