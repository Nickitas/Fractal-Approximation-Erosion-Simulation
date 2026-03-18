package cli

import "coastal-geometry/internal/domain/coastline"

func runAllCommand(app *App) error {
	invalid := false

	sanity := coastline.MainCalculation(app.Base, app.Dataset, app.DataSource)
	if sanity.Checked && !sanity.Valid {
		invalid = true
	}
	if err := writeCoastlineSVG(app.Base, app.Config.OutputPath, "coastline.svg", app.Config.Command); err != nil {
		return err
	}

	runParadoxCommand(app)
	runKochOrganicMetrics(app.Base, app.Config.Iterations, organicKochOptions(app))

	if err := writeOrganicKochSVGSeries(app.Base, app.Config.Iterations, app.Config.OutputPath, organicKochOptions(app), "koch_iter"); err != nil {
		return err
	}

	assessment, err := runDimensionMetrics(app.Base, app.Config.Iterations, organicKochOptions(app))
	if err != nil {
		return err
	}
	if !assessment.Valid {
		invalid = true
	}
	if invalid {
		printInvalidResult()
	}
	return nil
}
