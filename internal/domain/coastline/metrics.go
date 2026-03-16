package coastline

import (
	"fmt"
	"path/filepath"
	"strings"

	"coastal-geometry/internal/domain/geometry"
)

func MainCalculation(coast []geometry.LatLon, inputPath string) SanityCheckResult {
	segmentCount := 0
	if len(coast) > 1 {
		segmentCount = len(coast) - 1
	}

	fmt.Println(strings.Repeat("═", 80))
	fmt.Println("\tБЕРЕГОВАЯ ЛИНИЯ ЧЁРНОГО МОРЯ")
	fmt.Println(strings.Repeat("═", 80))

	fmt.Printf("\nКоличество точек:                        %d\n", len(coast))
	fmt.Printf("Количество сегментов:                    %d\n", segmentCount)

	totalLength := geometry.PolylineLength(coast)
	sanity := SanityCheck(filepath.Base(inputPath), totalLength)
	fmt.Printf("Общая длина береговой линии:              %.0f км\n", totalLength)
	if segmentCount > 0 {
		fmt.Printf("Средняя длина сегмента:                   %.1f км\n\n", totalLength/float64(segmentCount))
	} else {
		fmt.Printf("Средняя длина сегмента:                   0.0 км\n\n")
	}

	if sanity.Warning != "" {
		fmt.Println(sanity.Warning)
		fmt.Println()
	}

	fmt.Println("Ключевые точки береговой линии:")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("%-4s %-11s %-11s %-25s\n", "№", "Широта", "Долгота", "Город / ориентир")
	fmt.Println(strings.Repeat("─", 80))

	for i, p := range coast {
		name := getLocationName(p)
		fmt.Printf("%-4d %-11.4f %-11.4f %-25s\n", i+1, p.Lat, p.Lon, name)
	}

	fmt.Println(strings.Repeat("═", 80))
	fmt.Printf("Итого: %.0f км\n", totalLength)
	return sanity
}
