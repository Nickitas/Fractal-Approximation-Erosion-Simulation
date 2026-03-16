package cli

import (
	"coastal-geometry/internal/domain/generators/koch"
	"coastal-geometry/internal/domain/geometry"
)

func runKochOrganicCommand(app *App) error {
	opts := organicKochOptions(app)
	runKochOrganicMetrics(app.Base, app.Config.Iterations, opts)
	return writeOrganicKochSVGSeries(app.Base, app.Config.Iterations, app.Config.OutputPath, opts, "koch_iter")
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
