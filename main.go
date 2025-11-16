package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func main() {
	var p1x, p1y, p2x, p2y float64

	fmt.Println("Укажите координаты точек")
	fmt.Print("x1: ")
	fmt.Scan(&p1x)

	fmt.Print("y1: ")
	fmt.Scan(&p1y)

	fmt.Print("x2: ")
	fmt.Scan(&p2x)

	fmt.Print("y2: ")
	fmt.Scan(&p2x)

	p1 := Point{p1x, p1y}
	p2 := Point{p2x, p2y}

	dx := p1.X - p2.X
	dy := p1.Y - p2.Y

	d := math.Sqrt(dx*dx + dy*dy)
	fmt.Println(`Эвклидово расстояние между точками:`, d)
}
