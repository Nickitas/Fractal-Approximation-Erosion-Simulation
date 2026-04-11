package cli

import (
	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
)

func runKochOrganicCommand(app *App) error {
	opts := organicKochOptions(app)
	runKochOrganicMetrics(app.ModelBase, app.Config.Iterations, opts)
	if err := writeOrganicKochSVGSeries(app.Base, app.ModelBase, app.Config.Iterations, app.Config.OutputPath, opts, app.Config.ErosionStrength, "koch_iter", "koch-organic", false, newExportContext(app)); err != nil {
		return err
	}
	return writeOrganicKochSVGSeries(app.Base, app.ModelBase, app.Config.Iterations, app.Config.OutputPath, opts, app.Config.ErosionStrength, "dimension_iter", "dimension-organic", true, newExportContext(app))
}

func runKochOrganicMetrics(base []geometry.LatLon, iterations int, opts koch.OrganicOptions) {
	koch.DemonstrateOrganic(base, iterations, opts)
}

func organicKochOptions(app *App) koch.OrganicOptions {
	return koch.OrganicOptions{
		Seed:            app.Config.Seed,
		AngleJitterDeg:  app.Config.AngleJitter,
		HeightJitterPct: app.Config.HeightJitter,
	}
}
