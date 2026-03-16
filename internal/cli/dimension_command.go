package cli

import (
	"coastal-geometry/coastline"
	"coastal-geometry/fractal"
	"coastal-geometry/koch"
	"fmt"
	"math"
	"strings"
)

func runDimensionCommand(app *App) error {
	if err := writeKochSVGSeries(app.Base, app.Config.Iterations, app.Config.OutputPath); err != nil {
		return err
	}
	return runDimensionMetrics(app.Base, app.Config.Iterations)
}

func runDimensionMetrics(base []coastline.LatLon, maxIterations int) error {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\tВЫЧИСЛЕНИЕ ФРАКТАЛЬНОЙ РАЗМЕРНОСТИ (box-counting)")
	fmt.Println(strings.Repeat("=", 80))

	theoreticalD := math.Log(4) / math.Log(3)

	fmt.Printf("Теоретическая размерность кривой Коха: D = log(4)/log(3) ≈ %.5f\n\n", theoreticalD)

	fmt.Printf("%-5s %-10s %-12s %-12s %-12s\n", "Итер.", "Точек", "Длина, км", "D (расчёт)", "Δ от теории")
	fmt.Println(strings.Repeat("─", 80))

	for iter := 0; iter <= maxIterations; iter++ {
		curve := koch.KochCurve(base, iter)
		length := coastline.PolylineLength(curve)
		dimension := fractal.FractalDimension(curve)

		delta := ""
		if iter >= 2 {
			delta = fmt.Sprintf("%+.5f", dimension-theoreticalD)
		}

		fmt.Printf("%-5d %-10d %-12.0f %-12.5f %s\n", iter, len(curve), length, dimension, delta)
	}

	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("Чем выше итерация — тем точнее оценка размерности → 1.26186")
	fmt.Println("Это доказывает: наша кривая Коха — настоящий фрактал!\t")
	return nil
}
