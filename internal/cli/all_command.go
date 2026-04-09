package cli

import "coastal-geometry/internal/domain/coastline"

func runAllCommand(app *App) error {
	invalid := false

	sanity := coastline.MainCalculation(app.Base, app.Dataset, app.DataSource)
	if sanity.Checked && !sanity.Valid {
		invalid = true
	}
	if err := writeCoastlineSVG(app.Base, app.RenderBase, app.Config.OutputPath, "coastline.svg", newExportContext(app)); err != nil {
		return err
	}

	runParadoxCommand(app)
	runKochOrganicMetrics(app.ModelBase, app.Config.Iterations, organicKochOptions(app))

	if err := writeOrganicKochSVGSeries(app.Base, app.ModelBase, app.Config.Iterations, app.Config.OutputPath, organicKochOptions(app), app.Config.ErosionStrength, "koch_iter", "koch-organic", false, newExportContext(app)); err != nil {
		return err
	}

	assessment, err := runDimensionMetrics(app.ModelBase, app.Config.Iterations, organicKochOptions(app))
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
