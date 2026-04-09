package geometry

import (
	"math"
	"math/rand"
	"time"
)

// Erode applies a Gaussian-distributed random displacement to every point.
// strength is the standard deviation of the displacement in meters; zero or
// negative values return a clone of the input without changes.
func Erode(points []LatLon, strength float64) []LatLon {
	return erodeWithRand(points, strength, rand.New(rand.NewSource(time.Now().UnixNano())))
}

// ErodeWithSeed mirrors Erode but allows a fixed seed for reproducible output.
func ErodeWithSeed(points []LatLon, strength float64, seed int64) []LatLon {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return erodeWithRand(points, strength, rand.New(rand.NewSource(seed)))
}

// SimulateErosion runs multiple erosion steps and returns snapshot after each step,
// including the initial state at index 0.
func SimulateErosion(points []LatLon, steps int, strength float64) [][]LatLon {
	return SimulateErosionWithSeed(points, steps, strength, time.Now().UnixNano())
}

// SimulateErosionWithSeed is deterministic for a fixed seed.
func SimulateErosionWithSeed(points []LatLon, steps int, strength float64, seed int64) [][]LatLon {
	if steps < 0 {
		steps = 0
	}

	rng := rand.New(rand.NewSource(seed))
	snapshots := make([][]LatLon, steps+1)

	current := clonePoints(points)
	snapshots[0] = current

	for i := 1; i <= steps; i++ {
		current = erodeWithRand(current, strength, rng)
		snapshots[i] = current
	}
	return snapshots
}

func erodeWithRand(points []LatLon, strength float64, rng *rand.Rand) []LatLon {
	if len(points) == 0 {
		return nil
	}
	if strength <= 0 || rng == nil {
		return clonePoints(points)
	}

	// Use mean latitude to approximate meters-to-degrees conversion.
	refLat := 0.0
	for _, p := range points {
		refLat += p.Lat
	}
	refLat /= float64(len(points))

	metersPerDegLat := 111194.9
	metersPerDegLon := metersPerDegLat * math.Cos(refLat*math.Pi/180)
	if math.Abs(metersPerDegLon) < 1e-9 {
		metersPerDegLon = metersPerDegLat
	}

	eroded := make([]LatLon, len(points))
	firstShiftLat := 0.0
	firstShiftLon := 0.0
	closed := isClosedPolyline(points)

	for i, p := range points {
		dx := rng.NormFloat64() * strength
		dy := rng.NormFloat64() * strength

		if closed {
			if i == 0 {
				firstShiftLat = dy
				firstShiftLon = dx
			}
			if i == len(points)-1 {
				dy = firstShiftLat
				dx = firstShiftLon
			}
		}

		eroded[i] = LatLon{
			Lat: p.Lat + dy/metersPerDegLat,
			Lon: p.Lon + dx/metersPerDegLon,
		}
	}

	return eroded
}
