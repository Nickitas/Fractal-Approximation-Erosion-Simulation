// fractal/dimension.go
package fractal

import (
	"math"

	"coastal-geometry/coastline"
)

// Приближённо вычисляет фрактальную размерность методом box-counting
// Используем сетку из квадратов размера epsilon
// Возвращает D ≈ log(N)/log(1/ε)
func FractalDimension(points []coastline.LatLon) float64 {
	if len(points) < 10 {
		return 1.0
	}

	// Преобразуем все точки в метры (эквидистантная проекция)
	meters := make([]Point2D, len(points))
	for i, p := range points {
		meters[i] = latLonToMeters(p)
	}

	// Находим границы в метрах
	minX, maxX, minY, maxY := bboxMeters(meters)
	width := maxX - minX
	height := maxY - minY
	bboxSize := math.Max(width, height)
	if bboxSize < 1 {
		return 1.0
	}

	// Адаптивные масштабы: от 1/8 до 1/256 от размера объекта
	var logEps, logN []float64
	for factor := 8.0; factor <= 256; factor *= 2 {
		eps := bboxSize / factor
		n := boxesCoveredMeters(meters, eps, minX, minY)
		if n > 1 {
			logEps = append(logEps, math.Log(1.0/eps))
			logN = append(logN, math.Log(float64(n)))
		}
	}

	if len(logEps) < 3 {
		return 1.0
	}

	D := linearRegressionSlope(logEps, logN)
	if D < 0.5 || D > 3.0 { // защита от выбросов
		return 1.0
	}
	return D
}

// Point2D — координаты в метрах
type Point2D struct{ X, Y float64 }

// latLonToMeters — простая, но точная эквидистантная проекция для Чёрного моря
func latLonToMeters(p coastline.LatLon) Point2D {
	const (
		// Средняя широта Чёрного моря
		refLat = 43.5
		// 1° широты ≈ 111.2 км
		metersPerDegLat = 111194.9
		// 1° долготы на широте 43.5° ≈ 85 км
		metersPerDegLon = 87300.0
	)

	dLat := (p.Lat - refLat) * metersPerDegLat
	dLon := (p.Lon - 35.0) * metersPerDegLon // от условного центра

	return Point2D{X: dLon, Y: dLat}
}

func bboxMeters(points []Point2D) (minX, maxX, minY, maxY float64) {
	if len(points) == 0 {
		return 0, 0, 0, 0
	}
	minX, minY = points[0].X, points[0].Y
	maxX, maxY = points[0].X, points[0].Y
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return
}

func boxesCoveredMeters(points []Point2D, eps, minX, minY float64) int {
	covered := make(map[[2]int]bool)
	for _, p := range points {
		i := int(math.Floor((p.Y - minY) / eps))
		j := int(math.Floor((p.X - minX) / eps))
		covered[[2]int{i, j}] = true
	}
	return len(covered)
}

func linearRegressionSlope(x, y []float64) float64 {
	n := float64(len(x))
	var sumX, sumY, sumXY, sumX2 float64
	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}
	return (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
}
