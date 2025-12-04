package main

import (
	"coastal-geometry/coastline"
	"coastal-geometry/fractal"
	"coastal-geometry/koch"
	"coastal-geometry/paradox"
	"fmt"
	"math"
	"strings"
)

func main() {
	coastline.MainCalculation()
	paradox.Demonstrate()
	koch.Demonstrate()

	demonstrateKochWithDimension()
}

func demonstrateKochWithDimension() {
	base := coastline.LoadCoastlineData()
	// baseLength := coastline.PolylineLength()(base)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\tВЫЧИСЛЕНИЕ ФРАКТАЛЬНОЙ РАЗМЕРНОСТИ (box-counting)")
	fmt.Println(strings.Repeat("=", 80))

	theoreticalD := math.Log(4) / math.Log(3) // ≈1.261859

	fmt.Printf("Теоретическая размерность кривой Коха: D = log(4)/log(3) ≈ %.5f\n\n", theoreticalD)

	fmt.Printf("%-5s %-10s %-12s %-12s %-12s\n", "Итер.", "Точек", "Длина, км", "D (расчёт)", "Δ от теории")
	fmt.Println(strings.Repeat("─", 80))

	for iter := 0; iter <= 10; iter++ {
		curve := koch.KochCurve(base, iter)
		length := coastline.PolylineLength(curve)
		D := fractal.FractalDimension(curve)

		delta := ""
		if iter >= 2 { // на первых итерациях оценка неточная
			delta = fmt.Sprintf("%+.5f", D-theoreticalD)
		}

		fmt.Printf("%-5d %-10d %-12.0f %-12.5f %s\n", iter, len(curve), length, D, delta)
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("Чем выше итерация — тем точнее оценка размерности → 1.26186")
	fmt.Println("Это доказывает: наша кривая Коха — настоящий фрактал!\t")
}
