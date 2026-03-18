package cli

import (
	"coastal-geometry/internal/domain/coastline"
	"coastal-geometry/internal/domain/geometry"
)

type App struct {
	Config           config
	Base             []geometry.LatLon
	Validation       coastline.ValidationReport
	DataSource       string
	Dataset          string
	LoadNotes        []string
	SourceInspection *coastline.SourceInspection
}

func NewApp(cfg config) (*App, error) {
	app := &App{Config: cfg}

	if cfg.Command == cmdSource {
		inspection, err := coastline.InspectSource(coastline.InspectOptions{
			LocalPath:    cfg.InputPath,
			RemoteURL:    cfg.SourceURL,
			SnapshotPath: cfg.OutputPath,
			Refresh:      cfg.Refresh,
		})
		if err != nil {
			return nil, err
		}
		app.SourceInspection = &inspection
		app.DataSource = inspection.Source
		app.Dataset = inspection.DatasetName
		app.LoadNotes = inspection.LoadWarnings
		return app, nil
	}

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
