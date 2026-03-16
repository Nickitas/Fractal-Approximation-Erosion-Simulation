package fractal

import (
	"math"
	"testing"

	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
)

func TestAnalyzeBoxCountingStraightLine(t *testing.T) {
	line := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0, Lon: 0.05},
		{Lat: 0, Lon: 0.10},
		{Lat: 0, Lon: 0.15},
	}

	analysis := AnalyzeBoxCounting(line)
	if !analysis.Valid {
		t.Fatal("expected valid analysis for a straight line")
	}

	if math.Abs(analysis.Dimension-1.0) > 0.15 {
		t.Fatalf("expected dimension near 1.0, got %.5f", analysis.Dimension)
	}
}

func TestAnalyzeBoxCountingKochCurveNearTheory(t *testing.T) {
	base := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0, Lon: 0.2},
	}

	curve := koch.KochCurve(base, 5)
	analysis := AnalyzeBoxCounting(curve)
	theoretical := math.Log(4) / math.Log(3)

	if !analysis.Valid {
		t.Fatal("expected valid analysis for Koch curve")
	}

	if math.Abs(analysis.Dimension-theoretical) > 0.10 {
		t.Fatalf("expected dimension near %.5f, got %.5f", theoretical, analysis.Dimension)
	}
}

func TestAnalyzeBoxCountingProvidesScaleDiagnostics(t *testing.T) {
	base := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0.03, Lon: 0.1},
		{Lat: 0, Lon: 0.2},
	}

	analysis := AnalyzeBoxCounting(base)
	if len(analysis.Samples) < minScaleSamples {
		t.Fatalf("expected at least %d scale samples, got %d", minScaleSamples, len(analysis.Samples))
	}

	if len(analysis.LocalDimensions) == 0 {
		t.Fatal("expected local slope diagnostics")
	}
}
