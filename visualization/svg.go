package visualization

import (
	"coastal-geometry/coastline"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	canvasWidth  = 1200
	canvasHeight = 800
	padding      = 60.0
)

func DrawSVG(points []coastline.LatLon, filename, title string) error {
	if len(points) < 2 {
		return fmt.Errorf("need at least 2 points to draw svg")
	}

	minLat, maxLat, minLon, maxLon := bounds(points)
	lonSpan := maxLon - minLon
	latSpan := maxLat - minLat
	if lonSpan == 0 {
		lonSpan = 1
	}
	if latSpan == 0 {
		latSpan = 1
	}

	usableWidth := float64(canvasWidth) - 2*padding
	usableHeight := float64(canvasHeight) - 2*padding
	scale := math.Min(usableWidth/lonSpan, usableHeight/latSpan)

	var polyline strings.Builder
	for i, point := range points {
		x := padding + (point.Lon-minLon)*scale
		y := float64(canvasHeight) - padding - (point.Lat-minLat)*scale
		if i > 0 {
			polyline.WriteByte(' ')
		}
		polyline.WriteString(fmt.Sprintf("%.2f,%.2f", x, y))
	}

	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <rect width="100%%" height="100%%" fill="#f7f4ea"/>
  <g>
    <text x="60" y="48" font-family="Helvetica, Arial, sans-serif" font-size="28" fill="#16324f">%s</text>
    <text x="60" y="76" font-family="Helvetica, Arial, sans-serif" font-size="14" fill="#4f6d7a">points: %d</text>
    <polyline fill="none" stroke="#1f6f8b" stroke-width="3" stroke-linejoin="round" stroke-linecap="round" points="%s"/>
  </g>
</svg>
`, canvasWidth, canvasHeight, canvasWidth, canvasHeight, escapeText(title), len(points), polyline.String())

	if err := os.WriteFile(filename, []byte(svg), 0o644); err != nil {
		return fmt.Errorf("write svg %q: %w", filename, err)
	}

	return nil
}

func bounds(points []coastline.LatLon) (minLat, maxLat, minLon, maxLon float64) {
	minLat, maxLat = points[0].Lat, points[0].Lat
	minLon, maxLon = points[0].Lon, points[0].Lon

	for _, point := range points[1:] {
		if point.Lat < minLat {
			minLat = point.Lat
		}
		if point.Lat > maxLat {
			maxLat = point.Lat
		}
		if point.Lon < minLon {
			minLon = point.Lon
		}
		if point.Lon > maxLon {
			maxLon = point.Lon
		}
	}

	return minLat, maxLat, minLon, maxLon
}

func escapeText(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		`'`, "&apos;",
	)
	return replacer.Replace(value)
}
