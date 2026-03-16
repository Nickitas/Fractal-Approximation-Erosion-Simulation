package svg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"coastal-geometry/internal/domain/geometry"
)

func TestDrawDocumentIncludesLayersScaleAndLengths(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "fractal.svg")

	err := DrawDocument(Document{
		Title:    "Test Fractal",
		Subtitle: "Overlay",
		Layers: []Layer{
			{
				Label:       "Исходная полилиния",
				Points:      []geometry.LatLon{{Lat: 0, Lon: 0}, {Lat: 0, Lon: 1}},
				LengthKM:    100,
				Stroke:      "#000000",
				StrokeWidth: 2,
				Opacity:     1,
				DashArray:   "6 4",
			},
			{
				Label:       "Итерация 1",
				Points:      []geometry.LatLon{{Lat: 0, Lon: 0}, {Lat: 0.2, Lon: 0.5}, {Lat: 0, Lon: 1}},
				LengthKM:    140,
				Stroke:      "#ff0000",
				StrokeWidth: 3,
				Opacity:     1,
			},
		},
		Meta: []string{"Текущая длина: 140 км"},
	}, filename)
	if err != nil {
		t.Fatalf("DrawDocument returned error: %v", err)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read svg: %v", err)
	}

	svg := string(content)
	for _, expected := range []string{
		"Исходная полилиния",
		"Итерация 1",
		"Масштаб",
		"100 км",
		"140 км",
		"stroke-dasharray",
	} {
		if !strings.Contains(svg, expected) {
			t.Fatalf("expected SVG to contain %q", expected)
		}
	}
}
