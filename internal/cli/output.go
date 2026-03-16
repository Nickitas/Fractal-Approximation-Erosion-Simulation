package cli

import (
	"coastal-geometry/coastline"
	"coastal-geometry/koch"
	"coastal-geometry/visualization"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeCoastlineSVG(points []coastline.LatLon, output, defaultName, command string) error {
	filename, err := resolveOutputPath(output, defaultName, command)
	if err != nil {
		return err
	}

	if err := visualization.DrawSVG(points, filename, "Black Sea Coastline"); err != nil {
		return err
	}

	fmt.Printf("SVG saved to %s\n", filename)
	return nil
}

func writeKochSVGSeries(base []coastline.LatLon, iterations int, output string) error {
	outputDir, err := resolveSeriesOutputDir(output)
	if err != nil {
		return err
	}

	for iter := 0; iter <= iterations; iter++ {
		filename := filepath.Join(outputDir, fmt.Sprintf("koch_iter_%d.svg", iter))
		curve := koch.KochCurve(base, iter)
		if err := visualization.DrawSVG(curve, filename, fmt.Sprintf("Koch Curve Iteration %d", iter)); err != nil {
			return err
		}
		fmt.Printf("SVG saved to %s\n", filename)
	}
	return nil
}

func resolveOutputPath(output, defaultName, command string) (string, error) {
	if output == "" {
		output = defaultOutputDir
	}

	if strings.HasSuffix(strings.ToLower(output), ".svg") {
		if command == cmdAll {
			return "", fmt.Errorf("command %q generates multiple SVG files, so --output must be a directory", command)
		}

		dir := filepath.Dir(output)
		if dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return "", fmt.Errorf("create output directory %q: %w", dir, err)
			}
		}

		return filepath.Abs(output)
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		return "", fmt.Errorf("create output directory %q: %w", output, err)
	}

	return filepath.Abs(filepath.Join(output, defaultName))
}

func resolveSeriesOutputDir(output string) (string, error) {
	if output == "" {
		output = defaultOutputDir
	}

	if strings.HasSuffix(strings.ToLower(output), ".svg") {
		dir := filepath.Dir(output)
		if dir == "." {
			dir = defaultOutputDir
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", fmt.Errorf("create output directory %q: %w", dir, err)
		}
		return filepath.Abs(dir)
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		return "", fmt.Errorf("create output directory %q: %w", output, err)
	}

	return filepath.Abs(output)
}
