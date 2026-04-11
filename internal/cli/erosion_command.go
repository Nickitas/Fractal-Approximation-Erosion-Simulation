package cli

import (
	"coastal-geometry/internal/domain/geometry"
	"fmt"
	"strings"
)

func runErosionCommand(app *App) error {
	steps := app.Config.Steps
	strength := app.Config.ErosionStrength
	seed := app.Config.Seed

	snapshots := geometry.SimulateErosionWithSeed(app.ModelBase, steps, strength, seed)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("\tЭРОЗИЯ: МНОГОШАГОВАЯ СИМУЛЯЦИЯ")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Шаги=%d, σ=%.1f м, seed=%d\n\n", steps, strength, seed)
	fmt.Printf("%-6s %-10s %-12s %-14s\n", "Шаг", "Точек", "Длина, км", "Площадь, км²")
	fmt.Println(strings.Repeat("-", 56))

	for i, state := range snapshots {
		length := geometry.PolylineLength(state)
		area := geometry.Area(state)
		fmt.Printf("%-6d %-10d %-12.0f %-14.0f\n", i, len(state), length, area)
	}

	return writeErosionSVGSeries(app.Base, app.ModelBase, snapshots, steps, strength, seed, app.Config.OutputPath, newExportContext(app))
}
