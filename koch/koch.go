package koch

import (
	"coastal-geometry/coastline"
	"fmt"
	"math"
	"strings"
)

// Максимальная глубина рекурсии
// 5 итераций = 4^5 = 1024× больше точек → ~15 000 точек нормально
// 7 итераций = 262 144× → ~4 млн точек → уже много
const MaxIterations = 10

// Применение итерацию кривой Коха к полилинии рекурсивно
func KochCurve(base []coastline.LatLon, iterations int) []coastline.LatLon {
	if iterations < 0 {
		iterations = 0
	}
	if iterations > MaxIterations {
		fmt.Printf("Предупреждение: слишком много итераций (%d). Ограничено до %d\n", iterations, MaxIterations)
		iterations = MaxIterations
	}

	if iterations == 0 {
		result := make([]coastline.LatLon, len(base))
		copy(result, base)
		return result
	}

	return kochRecursive(base, iterations)
}

func kochRecursive(points []coastline.LatLon, depth int) []coastline.LatLon {
	if depth == 1 {
		return kochIteration(points)
	}
	// depth > 1 → сначала строим предыдущую итерацию, потом ещё одну
	return kochIteration(kochRecursive(points, depth-1))
}

// Итерация коха: для каждого сегмента AB добавляем 4 новые точки → 5 точек вместо 2
func kochIteration(points []coastline.LatLon) []coastline.LatLon {
	if len(points) < 2 {
		return points
	}

	newPoints := make([]coastline.LatLon, 0, len(points)*4)

	for i := 0; i < len(points)-1; i++ {
		segment := kochSegment(points[i], points[i+1])
		newPoints = append(newPoints, segment...)
	}
	newPoints = append(newPoints, points[len(points)-1])

	return newPoints
}

// Построение 5и точек Коха на отрезке от a до b
func kochSegment(a, b coastline.LatLon) []coastline.LatLon {
	// Векторы в градусах
	vx := b.Lon - a.Lon
	vy := b.Lat - a.Lat

	// Треть длины сегмента
	thirdX := vx / 3.0
	thirdY := vy / 3.0

	// Точки деления
	p1 := coastline.LatLon{Lat: a.Lat + thirdY, Lon: a.Lon + thirdX}
	p3 := coastline.LatLon{Lat: a.Lat + 2*thirdY, Lon: a.Lon + 2*thirdX}

	// p2 — вершина равностороннего треугольника
	// Поворачиваем вектор от p1 к p3 на 60° против часовой стрелки
	dx := thirdX
	dy := thirdY
	// Поворот на +60°: (x,y) → (x*cos60 - y*sin60, x*sin60 + y*cos60)
	cos60 := 0.5
	sin60 := math.Sqrt(3) / 2
	p2x := dx*cos60 - dy*sin60
	p2y := dx*sin60 + dy*cos60

	p2 := coastline.LatLon{
		Lat: p1.Lat + p2y,
		Lon: p1.Lon + p2x,
	}

	return []coastline.LatLon{a, p1, p2, p3}
}

func Demonstrate() {
	base := coastline.LoadCoastlineData()
	baseLength := coastline.PolylineLength(base)

	fmt.Println(strings.Repeat("═", 80))
	fmt.Println("\tФРАКТАЛЬНАЯ БЕРЕГОВАЯ ЛИНИЯ ЧЁРНОГО МОРЯ — КРИВАЯ КОХА (рекурсивная)")
	fmt.Println(strings.Repeat("═", 90))

	fmt.Printf("Исходная полилиния: %d точек, длина = %.0f км\n\n", len(base), baseLength)

	fmt.Printf("%-5s %-10s %-15s %-15s %-12s\n", "Итер.", "Точек", "Длина, км", "Прирост", "× от исходной")
	fmt.Println(strings.Repeat("─", 80))

	prevLength := baseLength
	// prevPoints := len(base)

	for iter := 0; iter <= MaxIterations; iter++ {
		curve := KochCurve(base, iter)
		length := coastline.PolylineLength(curve)
		pointsCount := len(curve)

		growth := ""
		multiplier := ""
		if iter > 0 {
			growth = fmt.Sprintf("+%.0f км", length-prevLength)
			multiplier = fmt.Sprintf("%.3f×", length/baseLength)
		} else {
			multiplier = "1.000×"
		}

		theoretical := baseLength * math.Pow(4.0/3.0, float64(iter))

		fmt.Printf("%-5d %-10d %-15.0f %-15s %-12s (теория: %.0f км)\n",
			iter, pointsCount, length, growth, multiplier, theoretical)

		prevLength = length
		// prevPoints = pointsCount
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("Математическая формула: Lₙ = L₀ × (4/3)ⁿ\n")
	fmt.Printf("Фрактальная размерность D = log(4)/log(3) ≈ %.5f\n", math.Log(4)/math.Log(3))
	fmt.Printf("При n→∞ длина → ∞, но кривая остаётся в ограниченной области\n")
}
