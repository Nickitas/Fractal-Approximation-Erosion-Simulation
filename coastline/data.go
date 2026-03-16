package coastline

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultCoastlineJSONPath = "data/black-sea.json"

func LoadCoastlineData() ([]LatLon, error) {
	return LoadFromJSON(DefaultCoastlineJSONPath)
}

func LoadFromJSON(filename string) ([]LatLon, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read coastline json %q: %w", filename, err)
	}

	var points []LatLon
	if err := json.Unmarshal(data, &points); err != nil {
		return nil, fmt.Errorf("parse coastline json %q: %w", filename, err)
	}

	if len(points) < 2 {
		return nil, fmt.Errorf("coastline json %q must contain at least 2 points", filename)
	}

	for i, point := range points {
		if point.Lat < -90 || point.Lat > 90 {
			return nil, fmt.Errorf("coastline json %q has invalid latitude at index %d: %f", filename, i, point.Lat)
		}
		if point.Lon < -180 || point.Lon > 180 {
			return nil, fmt.Errorf("coastline json %q has invalid longitude at index %d: %f", filename, i, point.Lon)
		}
	}

	return points, nil
}
