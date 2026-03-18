package cli

import "testing"

func TestModelBaseTargetPointsShrinksWithIterations(t *testing.T) {
	lowIteration := modelBaseTargetPoints(1)
	highIteration := modelBaseTargetPoints(5)

	if lowIteration <= highIteration {
		t.Fatalf("expected higher iterations to lower the point budget, got iter1=%d iter5=%d", lowIteration, highIteration)
	}
	if highIteration < 4 {
		t.Fatalf("expected model base target to keep a minimal polyline budget, got %d", highIteration)
	}
}

func TestCommandUsesModelBase(t *testing.T) {
	if !commandUsesModelBase(cmdKoch) {
		t.Fatal("expected koch command to use simplified model base")
	}
	if commandUsesModelBase(cmdCoastline) {
		t.Fatal("expected coastline command to avoid synthetic model-base simplification")
	}
}
