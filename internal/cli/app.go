package cli

import (
	"coastal-geometry/internal/domain/coastline"
	"coastal-geometry/internal/domain/geometry"
)

type App struct {
	Config     config
	Base       []geometry.LatLon
	Validation coastline.ValidationReport
}

func NewApp(cfg config) (*App, error) {
	app := &App{Config: cfg}

	if commandNeedsCoastline(cfg.Command) {
		base, report, err := coastline.LoadFromJSON(cfg.InputPath)
		if err != nil {
			return nil, err
		}
		app.Base = base
		app.Validation = report
	}

	return app, nil
}
