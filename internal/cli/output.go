package cli

import (
	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
	svgrender "coastal-geometry/internal/render/svg"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeCoastlineSVG(points []geometry.LatLon, output, defaultName, command string) error {
	filename, err := resolveOutputPath(output, defaultName, command)
	if err != nil {
		return err
	}

	if err := svgrender.DrawDocument(svgrender.Document{
		Title:    "Береговая линия",
		Subtitle: "Исходная географическая полилиния",
		Layers: []svgrender.Layer{
			{
				Label:       "Исходная полилиния",
				Points:      points,
				LengthKM:    geometry.PolylineLength(points),
				Stroke:      "#1f6f8b",
				StrokeWidth: 3.5,
				Opacity:     1,
			},
		},
		Meta: []string{
			fmt.Sprintf("Точек: %d", len(points)),
			fmt.Sprintf("Сегментов: %d", max(len(points)-1, 0)),
		},
	}, filename); err != nil {
		return err
	}

	fmt.Printf("SVG saved to %s\n", filename)
	return nil
}

func writeKochSVGSeries(base []geometry.LatLon, iterations int, output string) error {
	return writeKochLikeSVGSeries(base, iterations, output, "koch_iter", "Классическая кривая Коха", func(points []geometry.LatLon, iter int) []geometry.LatLon {
		return koch.KochCurve(points, iter)
	})
}

func writeOrganicKochSVGSeries(base []geometry.LatLon, iterations int, output string, opts koch.OrganicOptions, prefix string) error {
	return writeKochLikeSVGSeries(base, iterations, output, prefix, "Organic Koch", func(points []geometry.LatLon, iter int) []geometry.LatLon {
		return koch.OrganicKochCurve(points, iter, opts)
	})
}

func writeKochLikeSVGSeries(base []geometry.LatLon, iterations int, output, prefix, title string, builder func([]geometry.LatLon, int) []geometry.LatLon) error {
	outputDir, err := resolveSeriesOutputDir(output)
	if err != nil {
		return err
	}

	curves := make([][]geometry.LatLon, iterations+1)
	lengths := make([]float64, iterations+1)
	for iter := 0; iter <= iterations; iter++ {
		curves[iter] = builder(base, iter)
		lengths[iter] = geometry.PolylineLength(curves[iter])
	}

	for iter := 0; iter <= iterations; iter++ {
		filename := filepath.Join(outputDir, fmt.Sprintf("%s_%d.svg", prefix, iter))
		layers := makeFractalLayers(curves[:iter+1], lengths[:iter+1])
		meta := []string{
			fmt.Sprintf("Показано итераций: %d", iter+1),
			fmt.Sprintf("Текущая итерация: %d", iter),
			fmt.Sprintf("Текущая длина: %.0f км", lengths[iter]),
		}
		if err := svgrender.DrawDocument(svgrender.Document{
			Title:    fmt.Sprintf("%s — итерация %d", title, iter),
			Subtitle: "На схеме наложены исходная полилиния и все итерации до текущей",
			Layers:   layers,
			Meta:     meta,
		}, filename); err != nil {
			return err
		}
		fmt.Printf("SVG saved to %s\n", filename)
	}
	return nil
}

func makeFractalLayers(curves [][]geometry.LatLon, lengths []float64) []svgrender.Layer {
	palette := []string{
		"#7a8b99",
		"#2c7a7b",
		"#1f6f8b",
		"#c06c3f",
		"#8b3f5c",
		"#6f5f1f",
		"#3f6b4b",
	}

	layers := make([]svgrender.Layer, 0, len(curves))
	for i := range curves {
		label := fmt.Sprintf("Итерация %d", i)
		strokeWidth := 2.0
		opacity := 0.34 + float64(i)*0.08
		if opacity > 1 {
			opacity = 1
		}
		dashArray := ""
		if i == 0 {
			label = "Исходная полилиния"
			strokeWidth = 1.8
			opacity = 0.9
			dashArray = "8 6"
		}
		if i == len(curves)-1 {
			strokeWidth = 3.6
			opacity = 1
		}

		layers = append(layers, svgrender.Layer{
			Label:       label,
			Points:      curves[i],
			LengthKM:    lengths[i],
			Stroke:      palette[i%len(palette)],
			StrokeWidth: strokeWidth,
			Opacity:     opacity,
			DashArray:   dashArray,
		})
	}
	return layers
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
