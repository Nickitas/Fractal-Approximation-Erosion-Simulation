package paradox

import (
	"fmt"
	"strings"

	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
)

func Demonstrate(base []geometry.LatLon, maxIterations int) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("\tПАРАДОКС БЕРЕГОВОЙ ЛИНИИ")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Демонстрация использует изменение масштаба измерения и добавление новых")
	fmt.Println("геометрических деталей через кривую Коха. Простое деление сегментов без")
	fmt.Println("изменения формы здесь не используется.")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Printf("%-8s %-12s %-16s %-18s %-24s\n", "Уровень", "Точек", "Сегментов", "Средний шаг, км", "Длина, км")
	fmt.Println(strings.Repeat("-", 80))

	prevLength := 0.0
	for level := 0; level <= maxIterations; level++ {
		curve := koch.KochCurve(base, level)
		length := geometry.PolylineLength(curve)
		segments := max(len(curve)-1, 0)
		avgStep := 0.0
		if segments > 0 {
			avgStep = length / float64(segments)
		}

		growth := " | —"
		if level > 0 {
			growth = fmt.Sprintf(" | +%.0f км (%.3fx)", length-prevLength, length/prevLength)
		}

		fmt.Printf("%-8d %-12d %-16d %-18.2f %-24s\n", level, len(curve), segments, avgStep, fmt.Sprintf("%.0f%s", length, growth))
		prevLength = length
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Вывод: длина растёт при уменьшении шага измерения, потому что на каждом")
	fmt.Println("уровне появляются новые геометрические детали, а не только дополнительные точки.")
}
