package paradox

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"coastal-geometry/coastline"
)

func addMidpointPerturbation(a, b coastline.LatLon, deviation float64) coastline.LatLon {
	mid := coastline.LatLon{
		Lat: (a.Lat + b.Lat) / 2,
		Lon: (a.Lon + b.Lon) / 2,
	}

	dLat := b.Lat - a.Lat
	dLon := b.Lon - a.Lon
	length := math.Hypot(dLat, dLon)
	if length < 1e-9 {
		return mid
	}

	perpLat := -dLon / length
	perpLon := dLat / length

	offset := (rand.Float64() - 0.5) * deviation
	mid.Lat += perpLat * offset
	mid.Lon += perpLon * offset

	return mid
}

func refineCoastline(base []coastline.LatLon, levels, currentLevel int, deviation float64) []coastline.LatLon {
	if currentLevel >= levels {
		return base
	}

	result := make([]coastline.LatLon, 0, len(base)*2)
	result = append(result, base[0])

	for i := 1; i < len(base); i++ {
		mid := addMidpointPerturbation(base[i-1], base[i], deviation)
		result = append(result, mid, base[i])
	}

	return refineCoastline(result, levels, currentLevel+1, deviation*0.5)
}

func Demonstrate() {
	rand.Seed(time.Now().UnixNano())

	base := []coastline.LatLon{
		{46.48, 30.73}, // Одесса
		{44.62, 33.53}, // Севастополь
		{43.70, 39.75}, // Сочи
		{41.65, 41.63}, // Батуми
		{41.28, 31.42}, // Синоп
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("\tПАРАДОКС БЕРЕГОВОЙ ЛИНИИ (демонстрация)")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("%-8s %-12s %-15s %-12s\n", "Уровень", "Точек", "Длина, км", "Прирост")
	fmt.Println(strings.Repeat("-", 80))

	prevLength := 0.0
	for level := 0; level <= 6; level++ {
		coast := refineCoastline(base, level, 0, 0.15)

		length := coastline.PolylineLength(coast)
		growth := "—"
		if level > 0 {
			growth = fmt.Sprintf("+%.0f км (%.1fx)", length-prevLength, length/prevLength)
		}

		fmt.Printf("%-8d %-12d %-15.0f %s\n", level, len(coast), length, growth)
		prevLength = length
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Вывод: чем детальнее измеряем — тем длиннее берег.")
	fmt.Println("При бесконечной детализации длина → ∞ (фрактал)")
}
