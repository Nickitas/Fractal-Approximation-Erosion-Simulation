package cli

import "coastal-geometry/paradox"

func runParadoxCommand(_ *App) error {
	paradox.Demonstrate()
	return nil
}
