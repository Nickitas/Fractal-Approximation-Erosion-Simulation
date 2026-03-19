package cli

import (
	"coastal-geometry/internal/domain/coastline"
	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCoastlineSVGCreatesMetricsSidecar(t *testing.T) {
	dir := t.TempDir()
	points := []geometry.LatLon{
		{Lat: 46.48, Lon: 30.73},
		{Lat: 46.49, Lon: 30.74},
		{Lat: 0, Lon: 5},
	}
	renderPoints := []geometry.LatLon{
		{Lat: 46.48, Lon: 30.73},
		{Lat: 0, Lon: 5},
	}

	err := writeCoastlineSVG(points, renderPoints, dir, "coastline.svg", exportContext{
		Command: cmdCoastline,
		Dataset: "test.json",
		Source:  "unit-test",
		Validation: coastline.ValidationReport{
			Fixes:    []string{"normalized"},
			Warnings: []string{"long segment", "duplicate location"},
		},
	})
	if err != nil {
		t.Fatalf("writeCoastlineSVG returned error: %v", err)
	}

	metricsPath := filepath.Join(dir, "coastline.metrics.json")
	data, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatalf("read coastline metrics: %v", err)
	}

	var metrics coastlineArtifactMetrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		t.Fatalf("unmarshal coastline metrics: %v", err)
	}

	if metrics.Command != canonicalCommandPath(cmdCoastline) {
		t.Fatalf("expected command %q, got %q", canonicalCommandPath(cmdCoastline), metrics.Command)
	}
	if metrics.Real.PointsCount != len(points) {
		t.Fatalf("expected %d real points, got %d", len(points), metrics.Real.PointsCount)
	}
	if metrics.Render.PointsCount != len(renderPoints) {
		t.Fatalf("expected %d render points, got %d", len(renderPoints), metrics.Render.PointsCount)
	}
	if len(metrics.Validation.Summary) != 2 {
		t.Fatalf("expected 2 structured validation summaries, got %+v", metrics.Validation)
	}
	if len(metrics.Highlights.LongSegments) != 1 {
		t.Fatalf("expected 1 highlighted long segment in metrics, got %+v", metrics.Highlights)
	}
	if metrics.Highlights.LongSegments[0].StartIndex != 2 || metrics.Highlights.LongSegments[0].EndIndex != 3 {
		t.Fatalf("unexpected highlighted segment indices: %+v", metrics.Highlights.LongSegments[0])
	}
	if len(metrics.Validation.DuplicateLocations) != 1 {
		t.Fatalf("expected 1 duplicate location summary, got %+v", metrics.Validation.DuplicateLocations)
	}
	if metrics.Validation.DuplicateLocations[0].Name != "Одесса, Украина" {
		t.Fatalf("unexpected duplicate location summary: %+v", metrics.Validation.DuplicateLocations[0])
	}
	if len(metrics.Validation.Warnings) != 2 {
		t.Fatalf("expected validation warnings to be persisted, got %+v", metrics.Validation)
	}

	svgContent, err := os.ReadFile(filepath.Join(dir, "coastline.svg"))
	if err != nil {
		t.Fatalf("read coastline svg: %v", err)
	}
	svg := string(svgContent)
	for _, expected := range []string{"Предупреждения", "Длинные сегменты", "Контроль геометрии", "Автоисправления", "Сводка"} {
		if !strings.Contains(svg, expected) {
			t.Fatalf("expected coastline SVG to contain %q", expected)
		}
	}
}

func TestWriteKochSVGSeriesShowsReferenceAndModelBase(t *testing.T) {
	dir := t.TempDir()
	originalBase := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0.03, Lon: 0.10},
		{Lat: 0, Lon: 0.20},
	}
	modelBase := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0, Lon: 0.20},
	}

	err := writeKochSVGSeries(originalBase, modelBase, 1, dir, exportContext{
		Command: cmdKoch,
		Dataset: "test.json",
		Source:  "unit-test",
	})
	if err != nil {
		t.Fatalf("writeKochSVGSeries returned error: %v", err)
	}

	svgPath := filepath.Join(dir, "koch_iter_0.svg")
	content, err := os.ReadFile(svgPath)
	if err != nil {
		t.Fatalf("read series svg: %v", err)
	}
	svg := string(content)
	for _, expected := range []string{
		"Реальная загруженная",
		"полилиния (справочно)",
		"Упрощённая база модели",
		"(итерация 0)",
		"Длина по итерациям",
		"Контроль геометрии",
		"Повторы ориентиров",
		"Сводка",
	} {
		if !strings.Contains(svg, expected) {
			t.Fatalf("expected SVG to contain %q", expected)
		}
	}

	metricsPath := filepath.Join(dir, "koch.metrics.json")
	data, err := os.ReadFile(metricsPath)
	if err != nil {
		t.Fatalf("read series metrics: %v", err)
	}

	var metrics fractalSeriesArtifactMetrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		t.Fatalf("unmarshal series metrics: %v", err)
	}

	if !metrics.ModelSimplification.Applied {
		t.Fatalf("expected model simplification to be marked as applied, got %+v", metrics.ModelSimplification)
	}
	if len(metrics.Validation.Summary) != 2 {
		t.Fatalf("expected stable validation summary rows in series metrics, got %+v", metrics.Validation.Summary)
	}
	if len(metrics.Highlights.LongSegments) != 0 {
		t.Fatalf("expected stable empty highlights schema for clean series metrics, got %+v", metrics.Highlights)
	}
	if len(metrics.Iterations) != 2 {
		t.Fatalf("expected 2 iterations in metrics, got %d", len(metrics.Iterations))
	}
	if metrics.Iterations[0].Theory == nil {
		t.Fatal("expected theory diagnostics for Koch series")
	}
}

func TestWriteOrganicKochSVGSeriesPersistsDimensionMetrics(t *testing.T) {
	dir := t.TempDir()
	base := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0.03, Lon: 0.10},
		{Lat: 0, Lon: 0.20},
	}

	err := writeOrganicKochSVGSeries(base, base, 1, dir, koch.OrganicOptions{Seed: 7}, "dimension_iter", "dimension", true, exportContext{
		Command: cmdDimension,
		Dataset: "test.json",
		Source:  "unit-test",
	})
	if err != nil {
		t.Fatalf("writeOrganicKochSVGSeries returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "dimension.metrics.json"))
	if err != nil {
		t.Fatalf("read dimension metrics: %v", err)
	}

	var metrics fractalSeriesArtifactMetrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		t.Fatalf("unmarshal dimension metrics: %v", err)
	}

	foundDimension := false
	for _, iteration := range metrics.Iterations {
		if iteration.Dimension != nil {
			foundDimension = true
			break
		}
	}
	if !foundDimension {
		t.Fatal("expected at least one dimension summary in exported metrics")
	}

	svgContent, err := os.ReadFile(filepath.Join(dir, "dimension_iter_1.svg"))
	if err != nil {
		t.Fatalf("read dimension svg: %v", err)
	}
	svg := string(svgContent)
	for _, expected := range []string{"Размерность D", "Оценка", "Теория", "Сводка"} {
		if !strings.Contains(svg, expected) {
			t.Fatalf("expected dimension SVG to contain %q", expected)
		}
	}
}
