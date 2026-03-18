package cli

func executeCommand(app *App) error {
	switch app.Config.Command {
	case cmdSource:
		return runSourceCommand(app)
	case cmdAll:
		return runAllCommand(app)
	case cmdCoastline:
		return runCoastlineCommand(app)
	case cmdParadox:
		return runParadoxCommand(app)
	case cmdKoch:
		return runKochCommand(app)
	case cmdKochOrganic:
		return runKochOrganicCommand(app)
	case cmdDimension:
		return runDimensionCommand(app)
	default:
		return errUnsupportedCommand(app.Config.Command)
	}
}
