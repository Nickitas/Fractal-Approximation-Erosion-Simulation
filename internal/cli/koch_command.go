package cli

import (
	"coastal-geometry/coastline"
	"coastal-geometry/koch"
)

func runKochCommand(app *App) error {
	runKochMetrics(app.Base, app.Config.Iterations)
	return writeKochSVGSeries(app.Base, app.Config.Iterations, app.Config.OutputPath)
}

func runKochMetrics(base []coastline.LatLon, iterations int) {
	koch.Demonstrate(base, iterations)
}
