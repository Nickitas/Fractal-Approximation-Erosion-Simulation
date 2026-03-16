package cli

import "coastal-geometry/internal/domain/simulations/paradox"

func runParadoxCommand(app *App) error {
	paradox.Demonstrate(app.Base, app.Config.Iterations)
	return nil
}
