package cli

type commandUX struct {
	Mode        string
	Summary     string
	RuntimeNote string
}

func canonicalCommandPath(command string) string {
	switch command {
	case cmdSource:
		return cmdSource
	case cmdCoastline:
		return cmdReal + " " + cmdCoastline
	case cmdParadox:
		return cmdModel + " " + cmdParadox
	case cmdKoch:
		return cmdModel + " " + cmdKoch
	case cmdKochOrganic:
		return cmdModel + " " + cmdKochOrganic
	case cmdDimension:
		return cmdModel + " " + cmdDimension
	default:
		return command
	}
}

func legacyAlias(command string) string {
	switch command {
	case cmdCoastline, cmdParadox, cmdKoch, cmdKochOrganic, cmdDimension:
		return command
	default:
		return ""
	}
}

func getCommandUX(command string) commandUX {
	switch command {
	case cmdSource:
		return commandUX{
			Mode:        "dataset inspection",
			Summary:     "shows source metadata and saves a raw local snapshot of the selected dataset",
			RuntimeNote: "the command inspects the raw source payload and writes a snapshot copy without running coastline metrics or synthetic model stages",
		}
	case cmdCoastline:
		return commandUX{
			Mode:        "real-data analysis",
			Summary:     "reports geometry and geodesic metrics for the loaded coastline itself",
			RuntimeNote: "reported length and coastline.svg correspond to the loaded coastline without synthetic transformations",
		}
	case cmdParadox:
		return commandUX{
			Mode:        "synthetic demonstration",
			Summary:     "uses the loaded coastline only as a base polyline, then adds synthetic detail for the paradox demo",
			RuntimeNote: "iteration 0 is the loaded coastline; higher levels are synthetic refinements, not direct real-world measurements",
		}
	case cmdKoch:
		return commandUX{
			Mode:        "synthetic demonstration",
			Summary:     "uses the loaded coastline as a base polyline for the classic Koch model",
			RuntimeNote: "iteration 0 is the loaded coastline; Koch iterations are synthetic model curves derived from it",
		}
	case cmdKochOrganic:
		return commandUX{
			Mode:        "synthetic demonstration",
			Summary:     "uses the loaded coastline as a base polyline for an organic fractal model",
			RuntimeNote: "iteration 0 is the loaded coastline; later iterations are synthetic and tuned by jitter parameters",
		}
	case cmdDimension:
		return commandUX{
			Mode:        "synthetic demonstration",
			Summary:     "estimates box-counting dimension on synthetic organic iterations built from the loaded coastline",
			RuntimeNote: "dimension diagnostics apply to the generated organic model, not directly to the raw coastline geometry alone",
		}
	case cmdAll:
		return commandUX{
			Mode:        "mixed pipeline",
			Summary:     "starts with real-data coastline metrics, then runs synthetic paradox and fractal model stages",
			RuntimeNote: "the first stage reports the loaded coastline itself; subsequent stages are synthetic demonstrations derived from that base geometry",
		}
	default:
		return commandUX{}
	}
}
