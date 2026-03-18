package coastline

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetchCoastlineDataParsesGeoJSONPolygon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/geo+json")
		fmt.Fprint(w, `{
			"type": "FeatureCollection",
			"features": [
				{
					"type": "Feature",
					"geometry": {
						"type": "Polygon",
						"coordinates": [[
							[5.0, 5.0],
							[6.0, 6.0],
							[7.0, 5.0],
							[5.0, 5.0]
						]]
					}
				},
				{
					"type": "Feature",
					"geometry": {
						"type": "Polygon",
						"coordinates": [[
							[30.73, 46.48],
							[32.49, 45.33],
							[34.10, 44.94],
							[39.75, 43.70],
							[41.63, 41.65],
							[30.73, 46.48]
						]]
					}
				}
			]
		}`)
	}))
	defer server.Close()

	points, err := fetchCoastlineData(server.Client(), server.URL, DefaultBlackSeaBounds)
	if err != nil {
		t.Fatalf("fetchCoastlineData returned error: %v", err)
	}

	if len(points) != 6 {
		t.Fatalf("expected 6 polygon points inside Black Sea bounds, got %d", len(points))
	}
	if points[0].Lat != 46.48 || points[0].Lon != 30.73 {
		t.Fatalf("unexpected first point: %+v", points[0])
	}
	if points[len(points)-1].Lat != 46.48 || points[len(points)-1].Lon != 30.73 {
		t.Fatalf("expected closed ring to remain intact, got %+v", points[len(points)-1])
	}
}

func TestLoadUsesRemoteGeoJSONWhenAvailable(t *testing.T) {
	dir := t.TempDir()
	fallbackPath := filepath.Join(dir, "fallback.json")
	if err := os.WriteFile(fallbackPath, []byte(`[
		{"lat": 10.0, "lon": 10.0},
		{"lat": 11.0, "lon": 11.0}
	]`), 0o644); err != nil {
		t.Fatalf("write fallback json: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/geo+json")
		fmt.Fprint(w, `{
			"type": "Feature",
			"geometry": {
				"type": "LineString",
				"coordinates": [
					[30.73, 46.48],
					[32.49, 45.33],
					[34.10, 44.94],
					[39.75, 43.70]
				]
			}
		}`)
	}))
	defer server.Close()

	result, err := Load(LoadOptions{
		LocalPath:    fallbackPath,
		RemoteURL:    server.URL,
		RemoteBounds: DefaultBlackSeaBounds,
		HTTPClient:   server.Client(),
	})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.Source != server.URL {
		t.Fatalf("expected remote source %q, got %q", server.URL, result.Source)
	}
	if len(result.LoadWarnings) != 0 {
		t.Fatalf("expected no load warnings, got %+v", result.LoadWarnings)
	}
	if len(result.Points) != 4 {
		t.Fatalf("expected 4 remote points, got %d", len(result.Points))
	}
}

func TestLoadPreservesClosedPolygonRing(t *testing.T) {
	dir := t.TempDir()
	fallbackPath := filepath.Join(dir, "fallback.json")
	if err := os.WriteFile(fallbackPath, []byte(`[
		{"lat": 10.0, "lon": 10.0},
		{"lat": 11.0, "lon": 11.0}
	]`), 0o644); err != nil {
		t.Fatalf("write fallback json: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/geo+json")
		fmt.Fprint(w, `{
			"type": "Feature",
			"geometry": {
				"type": "Polygon",
				"coordinates": [[
					[30.73, 46.48],
					[32.49, 45.33],
					[34.10, 44.94],
					[39.75, 43.70],
					[30.73, 46.48]
				]]
			}
		}`)
	}))
	defer server.Close()

	result, err := Load(LoadOptions{
		LocalPath:  fallbackPath,
		RemoteURL:  server.URL,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if len(result.Points) != 5 {
		t.Fatalf("expected closed ring with 5 points, got %d", len(result.Points))
	}
	if result.Points[0] != result.Points[len(result.Points)-1] {
		t.Fatalf("expected normalized polygon ring to stay closed, got first=%+v last=%+v", result.Points[0], result.Points[len(result.Points)-1])
	}
}

func TestLoadFallsBackToLocalJSONWhenRemoteFails(t *testing.T) {
	dir := t.TempDir()
	fallbackPath := filepath.Join(dir, "fallback.json")
	if err := os.WriteFile(fallbackPath, []byte(`[
		{"lat": 46.48, "lon": 30.73},
		{"lat": 41.65, "lon": 41.63}
	]`), 0o644); err != nil {
		t.Fatalf("write fallback json: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "temporary failure", http.StatusBadGateway)
	}))
	defer server.Close()

	result, err := Load(LoadOptions{
		LocalPath:    fallbackPath,
		RemoteURL:    server.URL,
		RemoteBounds: DefaultBlackSeaBounds,
		HTTPClient:   server.Client(),
	})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if result.Source != fallbackPath {
		t.Fatalf("expected fallback source %q, got %q", fallbackPath, result.Source)
	}
	if len(result.Points) != 2 {
		t.Fatalf("expected 2 fallback points, got %d", len(result.Points))
	}
	if len(result.LoadWarnings) != 1 {
		t.Fatalf("expected one load warning, got %+v", result.LoadWarnings)
	}
	if !strings.Contains(result.LoadWarnings[0], "using local fallback") {
		t.Fatalf("unexpected load warning: %+v", result.LoadWarnings)
	}
}
