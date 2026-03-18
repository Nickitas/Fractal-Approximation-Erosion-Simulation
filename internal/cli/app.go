package cli

import (
	"coastal-geometry/internal/domain/coastline"
	"coastal-geometry/internal/domain/geometry"
)

type App struct {
	Config     config
	Base       []geometry.LatLon
	Validation coastline.ValidationReport
	DataSource string
	Dataset    string
	LoadNotes  []string
}

func NewApp(cfg config) (*App, error) {
	app := &App{Config: cfg}

	if commandNeedsCoastline(cfg.Command) {
		result, err := coastline.Load(coastline.LoadOptions{
			LocalPath: cfg.InputPath,
			RemoteURL: cfg.SourceURL,
			Refresh:   cfg.Refresh,
		})
		if err != nil {
			return nil, err
		}
		app.Base = result.Points
		app.Validation = result.Validation
		app.DataSource = result.Source
		app.Dataset = result.DatasetName
		app.LoadNotes = result.LoadWarnings
	}

	return app, nil
}
