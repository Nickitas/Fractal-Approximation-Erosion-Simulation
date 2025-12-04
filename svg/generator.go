package svg

import (
	"fmt"
	"os"

	"coastal-geometry/coastline"
	"coastal-geometry/fractal"
	"coastal-geometry/koch"
)

const (
	width  = 1200
	height = 800
	margin = 80
)

func GenerateAllIterations(base []coastline.LatLon, maxIter int, folder string) error {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}

	for iter := 0; iter <= maxIter; iter++ {
		curve := koch.KochCurve(base, iter)
		filename := fmt.Sprintf("%s/koch_%02d.svg", folder, iter)
		if err := GenerateSVG(curve, filename, iter); err != nil {
			return err
		}
		fmt.Printf("Сохранено: %s (итерация %d, точек: %d, длина: %.0f км)\n",
			filename, iter, len(curve), coastline.PolylineLength(curve))
	}
	return nil
}

func GenerateSVG(points []coastline.LatLon, filename string, iter int) error {
	minLat, maxLat, minLon, maxLon := bounds(points)
	latRange := maxLat - minLat
	lonRange := maxLon - minLon
	if latRange < 1e-9 || lonRange < 1e-9 {
		return fmt.Errorf("слишком маленький диапазон")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// ВСЁ ИСПРАВЛЕНО: убрал %% → теперь % работает!
	fmt.Fprintln(f, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(f, `<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">
  <rect width="100%%" height="100%%" fill="#0f172a"/>
  <g transform="translate(%d,%d)">`, width, height, width, height, margin, margin)

	// Путь
	fmt.Fprint(f, `<path d="M`)
	x := (points[0].Lon - minLon) / lonRange * float64(width-2*margin)
	y := (maxLat - points[0].Lat) / latRange * float64(height-2*margin)
	fmt.Fprintf(f, "%.2f,%.2f", x, y)

	for _, p := range points[1:] {
		x = (p.Lon - minLon) / lonRange * float64(width-2*margin)
		y = (maxLat - p.Lat) / latRange * float64(height-2*margin)
		fmt.Fprintf(f, " L%.2f,%.2f", x, y)
	}

	// Цвет по итерации
	color := []string{"#60a5fa", "#93c5fd", "#dbeefe", "#fde047", "#fbbf24", "#f97316", "#ef4444"}[min(iter, 6)]
	strokeWidth := 2.0
	if iter >= 4 {
		strokeWidth = 1.0
	}

	fmt.Fprintf(f, `" stroke="%s" stroke-width="%.1f" fill="none"/>`, color, strokeWidth)

	// Подписи
	length := coastline.PolylineLength(points)
	dim := "—"
	if iter >= 2 {
		dim = fmt.Sprintf("%.5f", fractal.FractalDimension(points))
	}

	title := fmt.Sprintf("Кривая Коха — итерация %d", iter)
	info := fmt.Sprintf("Точек: %d | Длина: %.0f км | D ≈ %s", len(points), length, dim)

	fmt.Fprintf(f, `<text x="%d" y="40" fill="white" font-size="32" font-weight="bold">%s</text>`, width/2-400, title)
	fmt.Fprintf(f, `<text x="%d" y="80" fill="#94a3b8" font-size="20">%s</text>`, width/2-400, info)
	fmt.Fprintln(f, `</g></svg>`)
	return nil
}

func bounds(points []coastline.LatLon) (minLat, maxLat, minLon, maxLon float64) {
	if len(points) == 0 {
		return 0, 0, 0, 0
	}
	minLat, maxLat = points[0].Lat, points[0].Lat
	minLon, maxLon = points[0].Lon, points[0].Lon
	for _, p := range points {
		if p.Lat < minLat { minLat = p.Lat }
		if p.Lat > maxLat { maxLat = p.Lat }
		if p.Lon < minLon { minLon = p.Lon }
		if p.Lon > maxLon { maxLon = p.Lon }
	}
	return
}

func min(a, b int) int {
	if a < b { return a }
	return b
}