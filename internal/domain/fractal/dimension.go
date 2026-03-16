package fractal

import (
	"math"

	"coastal-geometry/internal/domain/geometry"
)

const (
	minScaleSamples       = 4
	minStableLocalSlopes  = 3
	minRegressionRSquared = 0.98
	maxLocalSlopeSpread   = 0.18
)

var defaultScaleFactors = []float64{4, 8, 16, 32, 64, 128, 256, 512}

type Point2D struct{ X, Y float64 }

type BoxCountingSample struct {
	ScaleFactor   float64
	RelativeScale float64
	BoxSizeMeters float64
	BoxesCovered  int
	LogInvScale   float64
	LogBoxes      float64
}

type BoxCountingAnalysis struct {
	Dimension          float64
	RegressionRSquared float64
	StableAcrossScales bool
	StabilitySpread    float64
	Samples            []BoxCountingSample
	LocalDimensions    []float64
	Valid              bool
}

func FractalDimension(points []geometry.LatLon) float64 {
	analysis := AnalyzeBoxCounting(points)
	if !analysis.Valid {
		return 1.0
	}
	return analysis.Dimension
}

func AnalyzeBoxCounting(points []geometry.LatLon) BoxCountingAnalysis {
	if len(points) < 2 {
		return BoxCountingAnalysis{}
	}

	meters := make([]Point2D, len(points))
	for i, p := range points {
		meters[i] = latLonToMeters(p)
	}

	minX, maxX, minY, maxY := bboxMeters(meters)
	width := maxX - minX
	height := maxY - minY
	bboxSize := math.Max(width, height)
	if bboxSize < 1 {
		return BoxCountingAnalysis{}
	}

	samples := make([]BoxCountingSample, 0, len(defaultScaleFactors))
	logInvScale := make([]float64, 0, len(defaultScaleFactors))
	logBoxes := make([]float64, 0, len(defaultScaleFactors))
	for _, factor := range defaultScaleFactors {
		boxSize := bboxSize / factor
		if boxSize <= 0 {
			continue
		}
		boxes := boxesCoveredMeters(meters, boxSize, minX, minY)
		if boxes <= 1 {
			continue
		}

		relativeScale := boxSize / bboxSize
		sample := BoxCountingSample{
			ScaleFactor:   factor,
			RelativeScale: relativeScale,
			BoxSizeMeters: boxSize,
			BoxesCovered:  boxes,
			LogInvScale:   math.Log(1.0 / relativeScale),
			LogBoxes:      math.Log(float64(boxes)),
		}

		samples = append(samples, sample)
		logInvScale = append(logInvScale, sample.LogInvScale)
		logBoxes = append(logBoxes, sample.LogBoxes)
	}

	if len(samples) < minScaleSamples {
		return BoxCountingAnalysis{Samples: samples}
	}

	slope, intercept := linearRegression(logInvScale, logBoxes)
	rSquared := regressionRSquared(logInvScale, logBoxes, slope, intercept)
	localDimensions := localSlopeSeries(logInvScale, logBoxes)
	spread := valueSpread(localDimensions)
	stable := len(localDimensions) >= minStableLocalSlopes &&
		rSquared >= minRegressionRSquared &&
		spread <= maxLocalSlopeSpread

	if slope < 0.5 || slope > 3.0 {
		return BoxCountingAnalysis{
			Samples:            samples,
			LocalDimensions:    localDimensions,
			RegressionRSquared: rSquared,
			StabilitySpread:    spread,
		}
	}

	return BoxCountingAnalysis{
		Dimension:          slope,
		RegressionRSquared: rSquared,
		StableAcrossScales: stable,
		StabilitySpread:    spread,
		Samples:            samples,
		LocalDimensions:    localDimensions,
		Valid:              true,
	}
}

func latLonToMeters(p geometry.LatLon) Point2D {
	const (
		refLat          = 43.5
		metersPerDegLat = 111194.9
		metersPerDegLon = 87300.0
	)

	dLat := (p.Lat - refLat) * metersPerDegLat
	dLon := (p.Lon - 35.0) * metersPerDegLon

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

func boxesCoveredMeters(points []Point2D, boxSize, minX, minY float64) int {
	covered := make(map[[2]int]struct{})
	for i := 1; i < len(points); i++ {
		markSegmentBoxes(covered, points[i-1], points[i], boxSize, minX, minY)
	}
	return len(covered)
}

func markSegmentBoxes(covered map[[2]int]struct{}, a, b Point2D, boxSize, minX, minY float64) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	distance := math.Hypot(dx, dy)
	steps := 1
	if boxSize > 0 {
		steps = int(math.Ceil(distance/(boxSize/2))) + 1
	}
	if steps < 2 {
		steps = 2
	}

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := a.X + dx*t
		y := a.Y + dy*t
		row := int(math.Floor((y - minY) / boxSize))
		col := int(math.Floor((x - minX) / boxSize))
		covered[[2]int{row, col}] = struct{}{}
	}
}

func localSlopeSeries(x, y []float64) []float64 {
	if len(x) != len(y) || len(x) < 2 {
		return nil
	}

	slopes := make([]float64, 0, len(x)-1)
	for i := 1; i < len(x); i++ {
		denominator := x[i] - x[i-1]
		if math.Abs(denominator) < 1e-12 {
			continue
		}
		slopes = append(slopes, (y[i]-y[i-1])/denominator)
	}
	return slopes
}

func valueSpread(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	minValue := values[0]
	maxValue := values[0]
	for _, value := range values[1:] {
		if value < minValue {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue - minValue
}

func linearRegression(x, y []float64) (slope, intercept float64) {
	n := float64(len(x))
	var sumX, sumY, sumXY, sumX2 float64
	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}
	denominator := n*sumX2 - sumX*sumX
	if math.Abs(denominator) < 1e-12 {
		return 0, 0
	}
	slope = (n*sumXY - sumX*sumY) / denominator
	intercept = (sumY - slope*sumX) / n
	return slope, intercept
}

func regressionRSquared(x, y []float64, slope, intercept float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	var meanY float64
	for _, value := range y {
		meanY += value
	}
	meanY /= float64(len(y))

	var ssTot, ssRes float64
	for i := range x {
		predicted := slope*x[i] + intercept
		residual := y[i] - predicted
		total := y[i] - meanY
		ssRes += residual * residual
		ssTot += total * total
	}

	if ssTot < 1e-12 {
		return 1
	}
	return 1 - ssRes/ssTot
}
