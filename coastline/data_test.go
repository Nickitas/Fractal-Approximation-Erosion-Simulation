package coastline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromJSON(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "coast.json")
	content := `[
		{"lat": 46.48, "lon": 30.73},
		{"lat": 41.65, "lon": 41.63}
	]`

	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp json: %v", err)
	}

	points, err := LoadFromJSON(filename)
	if err != nil {
		t.Fatalf("LoadFromJSON returned error: %v", err)
	}

	if len(points) != 2 {
		t.Fatalf("expected 2 points, got %d", len(points))
	}

	if points[0].Lat != 46.48 || points[0].Lon != 30.73 {
		t.Fatalf("unexpected first point: %+v", points[0])
	}
}

func TestLoadFromJSONRejectsInvalidLatitude(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "invalid-coast.json")
	content := `[
		{"lat": 146.48, "lon": 30.73},
		{"lat": 41.65, "lon": 41.63}
	]`

	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp json: %v", err)
	}

	_, err := LoadFromJSON(filename)
	if err == nil {
		t.Fatal("expected error for invalid latitude, got nil")
	}

	if !strings.Contains(err.Error(), "invalid latitude") {
		t.Fatalf("unexpected error: %v", err)
	}
}
