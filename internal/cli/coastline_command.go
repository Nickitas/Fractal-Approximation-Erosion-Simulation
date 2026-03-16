package cli

import "coastal-geometry/coastline"

func runCoastlineCommand(app *App) error {
	coastline.MainCalculation(app.Base)
	return writeCoastlineSVG(app.Base, app.Config.OutputPath, "coastline.svg", app.Config.Command)
}
