package cli

import (
	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
)

func runKochCommand(app *App) error {
	report := runKochMetrics(app.ModelBase, app.Config.Iterations)
	if !report.Valid {
		printInvalidResult()
	}
	return writeKochSVGSeries(app.Base, app.ModelBase, app.Config.Iterations, app.Config.OutputPath, newExportContext(app))
}

func runKochMetrics(base []geometry.LatLon, iterations int) koch.TheoryCheckReport {
	return koch.Demonstrate(base, iterations)
}
