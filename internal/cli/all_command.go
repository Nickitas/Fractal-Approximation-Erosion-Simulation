package cli

import "coastal-geometry/coastline"

func runAllCommand(app *App) error {
	coastline.MainCalculation(app.Base)
	if err := writeCoastlineSVG(app.Base, app.Config.OutputPath, "coastline.svg", app.Config.Command); err != nil {
		return err
	}

	runParadoxCommand(app)
	runKochMetrics(app.Base, app.Config.Iterations)

	if err := writeKochSVGSeries(app.Base, app.Config.Iterations, app.Config.OutputPath); err != nil {
		return err
	}

	return runDimensionMetrics(app.Base, app.Config.Iterations)
}
