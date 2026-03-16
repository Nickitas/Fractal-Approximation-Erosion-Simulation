package coastline

import (
	"encoding/json"
	"fmt"
	"os"

	"coastal-geometry/internal/domain/geometry"
)

const DefaultCoastlineJSONPath = "data/black-sea.json"

type ValidationReport struct {
	Fixes    []string
	Warnings []string
}

func LoadFromJSON(filename string) ([]geometry.LatLon, ValidationReport, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, ValidationReport{}, fmt.Errorf("read coastline json %q: %w", filename, err)
	}

	var points []geometry.LatLon
	if err := json.Unmarshal(data, &points); err != nil {
		return nil, ValidationReport{}, fmt.Errorf("parse coastline json %q: %w", filename, err)
	}

	if len(points) < 2 {
		return nil, ValidationReport{}, fmt.Errorf("coastline json %q must contain at least 2 points", filename)
	}

	for i, point := range points {
		if point.Lat < -90 || point.Lat > 90 {
			return nil, ValidationReport{}, fmt.Errorf("coastline json %q has invalid latitude at index %d: %f", filename, i, point.Lat)
		}
		if point.Lon < -180 || point.Lon > 180 {
			return nil, ValidationReport{}, fmt.Errorf("coastline json %q has invalid longitude at index %d: %f", filename, i, point.Lon)
		}
	}

	normalized, report, err := validateAndNormalizePoints(points)
	if err != nil {
		return nil, ValidationReport{}, fmt.Errorf("validate coastline json %q: %w", filename, err)
	}

	return normalized, report, nil
}
