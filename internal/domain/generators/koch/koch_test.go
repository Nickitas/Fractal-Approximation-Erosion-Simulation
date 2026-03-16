package koch

import (
	"math"
	"testing"

	"coastal-geometry/internal/domain/geometry"
)

func TestTheoreticalLength(t *testing.T) {
	baseLength := 90.0
	got := TheoreticalLength(baseLength, 2)
	want := 160.0

	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("TheoreticalLength() = %v, want %v", got, want)
	}
}

func TestTheoryErrorPercent(t *testing.T) {
	got := TheoryErrorPercent(98, 100)
	want := 2.0

	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("TheoryErrorPercent() = %v, want %v", got, want)
	}
}

func TestKochCurveMatchesTheoryForSingleSegment(t *testing.T) {
	base := []geometry.LatLon{
		{Lat: 0, Lon: 0},
		{Lat: 0, Lon: 0.1},
	}

	baseLength := geometry.PolylineLength(base)
	curve := KochCurve(base, 1)
	measuredLength := geometry.PolylineLength(curve)
	theoreticalLength := TheoreticalLength(baseLength, 1)
	errorPct := TheoryErrorPercent(measuredLength, theoreticalLength)

	if errorPct > maxTheoryErrorPct {
		t.Fatalf("expected theory error <= %.2f%%, got %.4f%%", maxTheoryErrorPct, errorPct)
	}
}
