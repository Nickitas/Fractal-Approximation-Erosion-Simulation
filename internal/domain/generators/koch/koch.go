package koch

import (
	"coastal-geometry/internal/domain/geometry"
	"fmt"
	"math"
	"strings"
)

const MaxIterations = 10
const maxTheoryErrorPct = 2.0

type TheoryCheckSample struct {
	Iteration        int
	PointsCount      int
	MeasuredLengthKM float64
	TheoreticalKM    float64
	ErrorKM          float64
	ErrorPercent     float64
}

type TheoryCheckReport struct {
	Samples []TheoryCheckSample
	Valid   bool
}

func KochCurve(base []geometry.LatLon, iterations int) []geometry.LatLon {
	if iterations < 0 {
		iterations = 0
	}
	if iterations > MaxIterations {
		fmt.Printf("Предупреждение: слишком много итераций (%d). Ограничено до %d\n", iterations, MaxIterations)
		iterations = MaxIterations
	}

	if iterations == 0 {
		result := make([]geometry.LatLon, len(base))
		copy(result, base)
		return result
	}

	return kochRecursive(base, iterations)
}

func kochRecursive(points []geometry.LatLon, depth int) []geometry.LatLon {
	if depth == 1 {
		return kochIteration(points)
	}
	return kochIteration(kochRecursive(points, depth-1))
}

func kochIteration(points []geometry.LatLon) []geometry.LatLon {
	if len(points) < 2 {
		return points
	}

	newPoints := make([]geometry.LatLon, 0, len(points)*4)
	for i := 0; i < len(points)-1; i++ {
		segment := kochSegment(points[i], points[i+1])
		newPoints = append(newPoints, segment...)
	}
	newPoints = append(newPoints, points[len(points)-1])

	return newPoints
}

func kochSegment(a, b geometry.LatLon) []geometry.LatLon {
	vx := b.Lon - a.Lon
	vy := b.Lat - a.Lat

	thirdX := vx / 3.0
	thirdY := vy / 3.0

	p1 := geometry.LatLon{Lat: a.Lat + thirdY, Lon: a.Lon + thirdX}
	p3 := geometry.LatLon{Lat: a.Lat + 2*thirdY, Lon: a.Lon + 2*thirdX}

	dx := thirdX
	dy := thirdY
	cos60 := 0.5
	sin60 := math.Sqrt(3) / 2
	p2x := dx*cos60 - dy*sin60
	p2y := dx*sin60 + dy*cos60

	p2 := geometry.LatLon{
		Lat: p1.Lat + p2y,
		Lon: p1.Lon + p2x,
	}

	return []geometry.LatLon{a, p1, p2, p3}
}

func TheoreticalLength(baseLength float64, iterations int) float64 {
	return baseLength * math.Pow(4.0/3.0, float64(iterations))
}

func TheoryError(measuredLength, theoreticalLength float64) float64 {
	return math.Abs(measuredLength - theoreticalLength)
}

func TheoryErrorPercent(measuredLength, theoreticalLength float64) float64 {
	if theoreticalLength == 0 {
		return 0
	}
	return TheoryError(measuredLength, theoreticalLength) / theoreticalLength * 100
}

func CheckTheoryConsistency(base []geometry.LatLon, maxIterations int) TheoryCheckReport {
	baseLength := geometry.PolylineLength(base)
	report := TheoryCheckReport{
		Samples: make([]TheoryCheckSample, 0, maxIterations+1),
		Valid:   true,
	}

	for iter := 0; iter <= maxIterations; iter++ {
		curve := KochCurve(base, iter)
		measuredLength := geometry.PolylineLength(curve)
		theoreticalLength := TheoreticalLength(baseLength, iter)
		errorKM := TheoryError(measuredLength, theoreticalLength)
		errorPct := TheoryErrorPercent(measuredLength, theoreticalLength)

		report.Samples = append(report.Samples, TheoryCheckSample{
			Iteration:        iter,
			PointsCount:      len(curve),
			MeasuredLengthKM: measuredLength,
			TheoreticalKM:    theoreticalLength,
			ErrorKM:          errorKM,
			ErrorPercent:     errorPct,
		})

		if errorPct > maxTheoryErrorPct {
			report.Valid = false
		}
	}

	return report
}

func Demonstrate(base []geometry.LatLon, maxIterations int) TheoryCheckReport {
	baseLength := geometry.PolylineLength(base)
	report := CheckTheoryConsistency(base, maxIterations)

	fmt.Println(strings.Repeat("═", 80))
	fmt.Println("\tФРАКТАЛЬНАЯ БЕРЕГОВАЯ ЛИНИЯ ЧЁРНОГО МОРЯ — КРИВАЯ КОХА (рекурсивная)")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Printf("Исходная полилиния: %d точек, длина = %.0f км\n\n", len(base), baseLength)

	fmt.Printf("%-5s %-10s %-15s %-15s %-15s %-12s\n", "Итер.", "Точек", "Измерено, км", "Теория, км", "Ошибка, км", "Ошибка, %")
	fmt.Println(strings.Repeat("─", 96))

	for _, sample := range report.Samples {
		fmt.Printf("%-5d %-10d %-15.0f %-15.0f %-15.2f %-12.2f\n",
			sample.Iteration,
			sample.PointsCount,
			sample.MeasuredLengthKM,
			sample.TheoreticalKM,
			sample.ErrorKM,
			sample.ErrorPercent)

		if sample.ErrorPercent > maxTheoryErrorPct {
			fmt.Println("WARNING: Koch implementation inconsistent with theory")
		}
	}

	fmt.Println(strings.Repeat("─", 96))
	fmt.Printf("Математическая формула: Lₙ = L₀ × (4/3)ⁿ\n")
	fmt.Printf("error = |L_measured - L_theory|\n")
	fmt.Printf("Порог предупреждения: %.0f%%\n", maxTheoryErrorPct)
	fmt.Printf("Фрактальная размерность D = log(4)/log(3) ≈ %.5f\n", math.Log(4)/math.Log(3))
	fmt.Printf("При n→∞ длина → ∞, но кривая остаётся в ограниченной области\n")
	return report
}
