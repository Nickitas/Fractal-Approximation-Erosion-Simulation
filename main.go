package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func (p Point) Distance(to Point) float64 {
	dx := p.X - to.X
	dy := p.Y - to.X

	return math.Sqrt(dx*dx + dy*dy)
}

func main() {
	p1 := Point{0, 0}
	p2 := Point{1000, 0}

	length := p1.Distance(p2)
	fmt.Printf("Длина береговой линии: %.2f км\n", length)
}
