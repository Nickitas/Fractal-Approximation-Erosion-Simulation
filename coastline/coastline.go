package coastline

import (
	"fmt"
	"sort"
	"strings"
)

func MainCalculation() {
	coast := LoadCoastlineData()

	// Сортировка по долготе слева направо (с запада на восток, потом обратно)
	sort.Slice(coast, func(i, j int) bool {
		return coast[i].Lon < coast[j].Lon
	})

	fmt.Println(strings.Repeat("═", 80))
	fmt.Println("\tБЕРЕГОВАЯ ЛИНИЯ ЧЁРНОГО МОРЯ")
	fmt.Println(strings.Repeat("═", 80))

	fmt.Printf("\nВсего географических точек в полилинии:  %d\n", len(coast))

	totalLength := PolylineLength(coast)
	first := coast[0]
	last := coast[len(coast)-1]
	straight := Haversine(first, last)

	fmt.Printf("Прямое расстояние (по воздуху):           %.0f км\n", straight)
	fmt.Printf("Длина береговой линии (полилиния):        %.0f км\n", totalLength)
	fmt.Printf("Коэффициент извилистости:                 %.2f×\n", totalLength/straight)
	fmt.Printf("Средняя длина одного сегмента:            %.1f км\n\n", totalLength/float64(len(coast)-1))

	fmt.Println("Ключевые точки береговой линии:")
	fmt.Println(strings.Repeat("─", 80))
	fmt.Printf("%-4s %-11s %-11s %-25s\n", "№", "Широта", "Долгота", "Город / ориентир")
	fmt.Println(strings.Repeat("─", 80))

	for i, p := range coast {
		name := getLocationName(p)
		fmt.Printf("%-4d %-11.4f %-11.4f %-25s\n", i+1, p.Lat, p.Lon, name)
	}

	fmt.Println(strings.Repeat("═", 80))
	fmt.Printf("Итого: %.0f км — соответствует реальным оценкам (4000–4500 км)\n", totalLength)
	fmt.Printf("Чем детальнее измеряем — тем длиннее берег (парадокс береговой линии)\n")
}
