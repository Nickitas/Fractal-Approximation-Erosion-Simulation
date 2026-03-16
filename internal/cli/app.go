package cli

import "coastal-geometry/coastline"

type App struct {
	Config config
	Base   []coastline.LatLon
}

func NewApp(cfg config) (*App, error) {
	app := &App{Config: cfg}

	if commandNeedsCoastline(cfg.Command) {
		base, err := coastline.LoadFromJSON(cfg.InputPath)
		if err != nil {
			return nil, err
		}
		app.Base = base
	}

	return app, nil
}
