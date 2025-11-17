package test

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Point struct {
	X, Y float64
}

func (p Point) distance(to Point) float64 {
	dx := p.X - to.X
	dy := p.Y - to.Y

	return math.Sqrt(dx*dx + dy*dy)
}

// Вычисление общей длинны полилинии (сумма всех сегментов)
func polylineLength(points []Point) float64 {
	if len(points) < 2 {
		return 0
	}

	var total float64
	for i := 1; i < len(points); i++ {
		total += points[i-1].distance(points[i])
	}

	return total
}

// Генератор n случайных точек в заданных границах
// X — от minX до maxX (сортировка по X)
// Y — отклонение от "основной линии" (например, Y=0), с амплитудой maxDeviation
func generateRandomPoints(n int, minX, maxX, maxDeviation float64) []Point {
	rand.Seed(time.Now().UnixNano())

	points := make([]Point, n)

	for i := 0; i < n; i++ {
		x := minX + rand.Float64()*(maxX-minX)
		y := (rand.Float64() - 0.5) * 2 * maxDeviation
		points[i] = Point{X: x, Y: y}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].X < points[j].X
	})

	return points
}

func GenRandomPoints() {
	const (
		startX       = 0
		endX         = 1200
		numPoints    = 10
		maxDeviation = 120
	)

	randomPoints := generateRandomPoints(numPoints, startX, endX, maxDeviation)

	coastline := append([]Point{{X: startX, Y: 0}}, randomPoints...)
	coastline = append(coastline, Point{X: endX, Y: 0})

	sort.Slice(coastline, func(i, j int) bool {
		return coastline[i].X < coastline[j].X
	})

	fmt.Println("Сгенерированная случайная береговая линия (полилиния):")
	for i, p := range coastline {
		fmt.Printf("Точка %2d: (%.2f, %.2f)\n", i+1, p.X, p.Y)
	}

	length := polylineLength(coastline)
	straight := endX - startX

	fmt.Printf("\nПрямая расстояние:          %.2f усл. ед.\n", float64(straight))
	fmt.Printf("Длина извилистой береговой линии: %.2f усл. ед.\n", length)
	fmt.Printf("Удлинение берега:          %.2f (в %.2f раза длиннее)\n", length-float64(straight), length/float64(straight))
}
