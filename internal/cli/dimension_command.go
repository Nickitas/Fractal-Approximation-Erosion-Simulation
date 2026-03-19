package cli

import (
	"coastal-geometry/internal/domain/fractal"
	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
	"fmt"
	"math"
	"strings"
)

const (
	theoryConvergenceTolerance = 0.05
	iterationConvergenceDelta  = 0.03
	minConvergedIterations     = 3
)

type dimensionIterationResult struct {
	Iteration int
	Analysis  fractal.BoxCountingAnalysis
}

type dimensionAssessment struct {
	Valid bool
}

func runDimensionCommand(app *App) error {
	opts := organicKochOptions(app)
	if err := writeOrganicKochSVGSeries(app.Base, app.ModelBase, app.Config.Iterations, app.Config.OutputPath, opts, "dimension_iter", "dimension", true, newExportContext(app)); err != nil {
		return err
	}
	assessment, err := runDimensionMetrics(app.ModelBase, app.Config.Iterations, opts)
	if err != nil {
		return err
	}
	if !assessment.Valid {
		printInvalidResult()
	}
	return nil
}

func runDimensionMetrics(base []geometry.LatLon, maxIterations int, opts koch.OrganicOptions) (dimensionAssessment, error) {
	theoreticalDimension := math.Log(4) / math.Log(3)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\tЭМПИРИЧЕСКАЯ ФРАКТАЛЬНАЯ РАЗМЕРНОСТЬ (box-counting)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("Organic model: seed=%d, angle jitter=±%.1f°, height jitter=±%.0f%%\n\n",
		opts.Seed, opts.AngleJitterDeg, opts.HeightJitterPct*100)
	fmt.Printf("Теоретический ориентир классической кривой Коха: %.5f\n\n", theoreticalDimension)

	fmt.Printf("%-5s %-10s %-12s %-12s %-8s %-8s %-10s %-10s %-8s\n",
		"Итер.", "Точек", "Длина, км", "D", "Масш.", "R²", "Разброс", "Δ к пред.", "Стаб.")
	fmt.Println(strings.Repeat("─", 104))

	results := make([]dimensionIterationResult, 0, maxIterations+1)
	prevDimension := 0.0
	prevValid := false
	for iter := 0; iter <= maxIterations; iter++ {
		curve := koch.OrganicKochCurve(base, iter, opts)
		length := geometry.PolylineLength(curve)
		analysis := fractal.AnalyzeBoxCounting(curve)
		results = append(results, dimensionIterationResult{Iteration: iter, Analysis: analysis})

		delta := "—"
		if prevValid && analysis.Valid {
			delta = fmt.Sprintf("%+.5f", analysis.Dimension-prevDimension)
		}

		stable := "no"
		dimensionValue := "n/a"
		rSquared := "n/a"
		spread := "n/a"
		if analysis.Valid {
			dimensionValue = fmt.Sprintf("%.5f", analysis.Dimension)
			rSquared = fmt.Sprintf("%.4f", analysis.RegressionRSquared)
			spread = fmt.Sprintf("%.4f", analysis.StabilitySpread)
			if analysis.StableAcrossScales {
				stable = "yes"
			}
			prevDimension = analysis.Dimension
			prevValid = true
		} else {
			prevValid = false
		}

		fmt.Printf("%-5d %-10d %-12.0f %-12s %-8d %-8s %-10s %-10s %-8s\n",
			iter, len(curve), length, dimensionValue, len(analysis.Samples), rSquared, spread, delta, stable)
	}

	fmt.Println(strings.Repeat("─", 104))
	return printDimensionAssessment(results, theoreticalDimension), nil
}

func printDimensionAssessment(results []dimensionIterationResult, theoreticalDimension float64) dimensionAssessment {
	valid := make([]dimensionIterationResult, 0, len(results))
	for _, result := range results {
		if result.Analysis.Valid {
			valid = append(valid, result)
		}
	}

	if len(valid) < minConvergedIterations {
		fmt.Println("Недостаточно валидных масштабов для оценки сходимости.")
		fmt.Println("Текущие результаты не дают оснований утверждать, что модель согласуется с теоретическим значением.")
		return dimensionAssessment{Valid: false}
	}

	tail := valid[len(valid)-minConvergedIterations:]
	convergedAcrossIterations := true
	for i := 1; i < len(tail); i++ {
		if math.Abs(tail[i].Analysis.Dimension-tail[i-1].Analysis.Dimension) > iterationConvergenceDelta {
			convergedAcrossIterations = false
			break
		}
	}

	stableAcrossScales := true
	for _, result := range tail {
		if !result.Analysis.StableAcrossScales {
			stableAcrossScales = false
			break
		}
	}

	finalDimension := tail[len(tail)-1].Analysis.Dimension
	deltaTheory := math.Abs(finalDimension - theoreticalDimension)

	fmt.Printf("Последняя оценка: D=%.5f, |D-D_theory|=%.5f\n", finalDimension, deltaTheory)
	fmt.Printf("Сходимость по последним %d итерациям: %v\n", minConvergedIterations, yesNo(convergedAcrossIterations))
	fmt.Printf("Стабильность на нескольких масштабах: %v\n", yesNo(stableAcrossScales))

	if convergedAcrossIterations && stableAcrossScales && deltaTheory <= theoryConvergenceTolerance {
		fmt.Println("Эмпирическая оценка согласуется с теоретическим ориентиром классической кривой Коха.")
		fmt.Println("Для organic-модели это ориентир, а не строгое доказательство.")
		return dimensionAssessment{Valid: true}
	}

	fmt.Println("Результаты не сходятся достаточно надёжно к теоретическому значению.")
	fmt.Println("По этим данным нельзя утверждать, что модель подтверждает теорию.")
	return dimensionAssessment{Valid: false}
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
