package svg

import (
	"coastal-geometry/internal/domain/geometry"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	canvasWidth   = 1440
	canvasHeight  = 900
	padding       = 56.0
	headerHeight  = 108.0
	sidebarWidth  = 320.0
	scaleBarYGap  = 42.0
	defaultStroke = "#1f6f8b"
)

type Layer struct {
	Label       string
	Points      []geometry.LatLon
	LengthKM    float64
	Stroke      string
	StrokeWidth float64
	Opacity     float64
	DashArray   string
}

type Document struct {
	Title    string
	Subtitle string
	Layers   []Layer
	Meta     []string
}

func DrawSVG(points []geometry.LatLon, filename, title string) error {
	return DrawDocument(Document{
		Title: title,
		Layers: []Layer{
			{
				Label:       "Исходная полилиния",
				Points:      points,
				LengthKM:    geometry.PolylineLength(points),
				Stroke:      defaultStroke,
				StrokeWidth: 3,
				Opacity:     1,
			},
		},
	}, filename)
}

func DrawDocument(doc Document, filename string) error {
	if len(doc.Layers) == 0 {
		return fmt.Errorf("need at least 1 layer to draw svg")
	}

	allPoints := flattenLayers(doc.Layers)
	if len(allPoints) < 2 {
		return fmt.Errorf("need at least 2 points to draw svg")
	}

	minLat, maxLat, minLon, maxLon := bounds(allPoints)
	lonSpan := maxLon - minLon
	latSpan := maxLat - minLat
	if lonSpan == 0 {
		lonSpan = 1
	}
	if latSpan == 0 {
		latSpan = 1
	}

	plotWidth := float64(canvasWidth) - sidebarWidth - 2*padding
	plotHeight := float64(canvasHeight) - headerHeight - padding
	scale := math.Min(plotWidth/lonSpan, plotHeight/latSpan)
	contentWidth := lonSpan * scale
	contentHeight := latSpan * scale
	originX := padding + (plotWidth-contentWidth)/2
	originY := headerHeight + (plotHeight-contentHeight)/2

	var layers strings.Builder
	for _, layer := range doc.Layers {
		polyline := projectPolyline(layer.Points, minLat, minLon, originX, originY, contentHeight, scale)
		layers.WriteString(fmt.Sprintf(
			`    <polyline fill="none" stroke="%s" stroke-width="%.2f" stroke-opacity="%.2f" stroke-linejoin="round" stroke-linecap="round"%s points="%s"/>`+"\n",
			escapeText(layerStroke(layer)),
			layerWidth(layer),
			layerOpacity(layer),
			layerDashAttribute(layer),
			polyline,
		))
	}

	sidebarX := padding + plotWidth + 28
	var legend strings.Builder
	legend.WriteString(fmt.Sprintf(
		`    <text x="%.0f" y="146" font-family="Helvetica, Arial, sans-serif" font-size="18" font-weight="700" fill="#16324f">Слои и длины</text>`+"\n",
		sidebarX,
	))
	for i, layer := range doc.Layers {
		y := 184.0 + float64(i)*34
		legend.WriteString(fmt.Sprintf(
			`    <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="%s" stroke-width="%.2f" stroke-opacity="%.2f"%s/>`+"\n",
			sidebarX,
			y,
			sidebarX+34,
			y,
			escapeText(layerStroke(layer)),
			layerWidth(layer),
			layerOpacity(layer),
			layerDashAttribute(layer),
		))
		legend.WriteString(fmt.Sprintf(
			`    <text x="%.0f" y="%.0f" font-family="Helvetica, Arial, sans-serif" font-size="14" fill="#16324f">%s</text>`+"\n",
			sidebarX+46,
			y+4,
			escapeText(fmt.Sprintf("%s — %.0f км", layer.Label, layer.LengthKM)),
		))
	}

	var meta strings.Builder
	for i, line := range doc.Meta {
		y := 608.0 + float64(i)*22
		meta.WriteString(fmt.Sprintf(
			`    <text x="%.0f" y="%.0f" font-family="Helvetica, Arial, sans-serif" font-size="13" fill="#4f6d7a">%s</text>`+"\n",
			sidebarX,
			y,
			escapeText(line),
		))
	}

	scaleBar := buildScaleBar(minLat, maxLat, minLon, maxLon, plotWidth, scale, padding, float64(canvasHeight)-padding-scaleBarYGap)

	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">
  <rect width="100%%" height="100%%" fill="#f7f4ea"/>
  <rect x="20" y="20" width="%d" height="%d" rx="28" fill="#fcfbf7" stroke="#d6d0c4"/>
  <rect x="%.0f" y="20" width="%.0f" height="%d" rx="24" fill="#f0ece2" stroke="#d6d0c4"/>
  <text x="56" y="58" font-family="Helvetica, Arial, sans-serif" font-size="30" font-weight="700" fill="#16324f">%s</text>
  <text x="56" y="86" font-family="Helvetica, Arial, sans-serif" font-size="14" fill="#4f6d7a">%s</text>
  <text x="56" y="112" font-family="Helvetica, Arial, sans-serif" font-size="13" fill="#6b7a87">SVG содержит исходную полилинию, фрактальные итерации, масштаб и длину по слоям.</text>
  <g>
%s  </g>
  <g>
%s  </g>
  <g>
%s  </g>
  <g>
%s  </g>
</svg>
`, canvasWidth, canvasHeight, canvasWidth, canvasHeight,
		canvasWidth-40, canvasHeight-40,
		padding+plotWidth+8, sidebarWidth-16, canvasHeight-40,
		escapeText(doc.Title),
		escapeText(doc.Subtitle),
		layers.String(),
		legend.String(),
		meta.String(),
		scaleBar,
	)

	if err := os.WriteFile(filename, []byte(svg), 0o644); err != nil {
		return fmt.Errorf("write svg %q: %w", filename, err)
	}

	return nil
}

func flattenLayers(layers []Layer) []geometry.LatLon {
	total := 0
	for _, layer := range layers {
		total += len(layer.Points)
	}

	points := make([]geometry.LatLon, 0, total)
	for _, layer := range layers {
		points = append(points, layer.Points...)
	}
	return points
}

func projectPolyline(points []geometry.LatLon, minLat, minLon, originX, originY, contentHeight, scale float64) string {
	var polyline strings.Builder
	for i, point := range points {
		x := originX + (point.Lon-minLon)*scale
		y := originY + contentHeight - (point.Lat-minLat)*scale
		if i > 0 {
			polyline.WriteByte(' ')
		}
		polyline.WriteString(fmt.Sprintf("%.2f,%.2f", x, y))
	}
	return polyline.String()
}

func buildScaleBar(minLat, maxLat, minLon, maxLon, plotWidth, scale, x, y float64) string {
	centerLat := (minLat + maxLat) / 2
	centerLon := (minLon + maxLon) / 2
	kmPerLonDegree := geometry.Haversine(
		geometry.LatLon{Lat: centerLat, Lon: centerLon},
		geometry.LatLon{Lat: centerLat, Lon: centerLon + 1},
	)
	if kmPerLonDegree <= 0 || scale <= 0 {
		return ""
	}

	kmPerPixel := kmPerLonDegree / scale
	targetKM := kmPerPixel * plotWidth * 0.22
	scaleKM := niceScaleLength(targetKM)
	barPixels := scaleKM / kmPerPixel

	return fmt.Sprintf(
		`    <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#16324f" stroke-width="3"/>
    <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#16324f" stroke-width="6"/>
    <line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#16324f" stroke-width="6"/>
    <text x="%.0f" y="%.0f" font-family="Helvetica, Arial, sans-serif" font-size="13" fill="#16324f">Масштаб ≈ %.0f км</text>`,
		x, y, x+barPixels, y,
		x, y-7, x, y+7,
		x+barPixels, y-7, x+barPixels, y+7,
		x, y-14, scaleKM,
	)
}

func niceScaleLength(value float64) float64 {
	if value <= 0 {
		return 1
	}

	power := math.Pow(10, math.Floor(math.Log10(value)))
	normalized := value / power

	switch {
	case normalized <= 1:
		return 1 * power
	case normalized <= 2:
		return 2 * power
	case normalized <= 5:
		return 5 * power
	default:
		return 10 * power
	}
}

func layerStroke(layer Layer) string {
	if layer.Stroke == "" {
		return defaultStroke
	}
	return layer.Stroke
}

func layerOpacity(layer Layer) float64 {
	if layer.Opacity <= 0 {
		return 1
	}
	return layer.Opacity
}

func layerWidth(layer Layer) float64 {
	if layer.StrokeWidth <= 0 {
		return 2
	}
	return layer.StrokeWidth
}

func layerDashAttribute(layer Layer) string {
	if layer.DashArray == "" {
		return ""
	}
	return fmt.Sprintf(` stroke-dasharray="%s"`, escapeText(layer.DashArray))
}

func bounds(points []geometry.LatLon) (minLat, maxLat, minLon, maxLon float64) {
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
