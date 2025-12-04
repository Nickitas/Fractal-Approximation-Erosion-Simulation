// wasm/main.go
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"

	"coastal-geometry/coastline"
	"coastal-geometry/fractal"
	"coastal-geometry/koch"
)

var currentIter = 0
var maxIter = 6
var base []coastline.LatLon

func main() {
	base = coastline.LoadCoastlineData()

	// Устанавливаем тёмный фон сразу
	doc := js.Global().Get("document")
	html := doc.Get("documentElement")
	body := doc.Get("body")
	html.Get("style").Set("height", "100vh")
	html.Get("style").Set("margin", "0")
	body.Get("style").Set("background", "#0f172a")
	body.Get("style").Set("fontFamily", "Arial, sans-serif")
	body.Get("style").Set("color", "#e2e8f0")
	body.Get("style").Set("textAlign", "center")
	body.Get("style").Set("padding", "20px")

	// Заголовок
	h1 := doc.Call("createElement", "h1")
	h1.Set("textContent", "Фрактальная береговая линия Чёрного моря")
	body.Call("appendChild", h1)

	// SVG
	svg := doc.Call("createElement", "svg")
	svg.Set("width", "1200")
	svg.Set("height", "800")
	svg.Set("viewBox", "0 0 1200 800")
	svg.Get("style").Set("background", "#0f172a")
	svg.Get("style").Set("border", "3px solid #475569")
	svg.Get("style").Set("borderRadius", "12px")
	svg.Set("id", "koch")
	body.Call("appendChild", svg)

	// Информация
	info := doc.Call("createElement", "div")
	info.Set("id", "info")
	info.Get("style").Set("margin", "20px")
	info.Get("style").Set("fontSize", "22px")
	body.Call("appendChild", info)

	// Кнопки
	btnPrev := createButton("Предыдущая итерация", func() { if currentIter > 0 { currentIter--; render() } })
	btnNext := createButton("Следующая итерация", func() { if currentIter < maxIter { currentIter++; render() } })
	btnAuto := createButton("Автоанимация", autoAnimate)

	body.Call("appendChild", btnPrev)
	body.Call("appendChild", btnNext)
	body.Call("appendChild", btnAuto)

	render() // первый кадр

	// Держим WASM живым
	select {}
}

func createButton(text string, fn func()) js.Value {
	doc := js.Global().Get("document")
	btn := doc.Call("createElement", "button")
	btn.Set("textContent", text)
	btn.Get("style").Set("padding", "12px 24px")
	btn.Get("style").Set("margin", "0 10px")
	btn.Get("style").Set("fontSize", "18px")
	btn.Get("style").Set("cursor", "pointer")
	btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	}))
	return btn
}

func render() {
	curve := koch.KochCurve(base, currentIter)
	length := coastline.PolylineLength(curve)
	d := "—"
	if currentIter >= 2 {
		d = fmt.Sprintf("%.5f", fractal.FractalDimension(curve))
	}

	svg := js.Global().Get("document").Call("getElementById", "koch")
	svg.Set("innerHTML", generateSVG(curve, currentIter))

	info := js.Global().Get("document").Call("getElementById", "info")
	info.Set("innerHTML", fmt.Sprintf(
		"<strong>Итерация %d</strong> → Точек: %d | Длина: %.0f км | D ≈ %s",
		currentIter, len(curve), length, d,
	))
}

func autoAnimate() {
	go func() {
		for i := 0; i <= maxIter; i++ {
			currentIter = i
			render()
			js.Global().Call("setTimeout", js.FuncOf(func(js.Value, []js.Value) any {
				return nil
			}), 1200)
		}
	}()
}

func generateSVG(points []coastline.LatLon, iter int) string {
	if len(points) == 0 { return "" }

	minLat, maxLat, minLon, maxLon := bounds(points)
	latR := maxLat - minLat
	lonR := maxLon - minLon

	colors := []string{"#60a5fa", "#93c5fd", "#dbeefe", "#fde047", "#fbbf24", "#f97316", "#ef4444"}
	color := colors[min(iter, len(colors)-1)]

	path := "M"
	x := (points[0].Lon-minLon)/lonR*1040 + 80
	y := (maxLat-points[0].Lat)/latR*640 + 80
	path += fmt.Sprintf("%.1f,%.1f", x, y)

	for _, p := range points[1:] {
		x = (p.Lon-minLon)/lonR*1040 + 80
		y = (maxLat-p.Lat)/latR*640 + 80
		path += fmt.Sprintf(" L%.1f,%.1f", x, y)
	}

	return fmt.Sprintf(`
    <rect width="100%%" height="100%%" fill="#0f172a"/>
    <path d="%s" stroke="%s" stroke-width="2" fill="none"/>
    <text x="80" y="80" fill="white" font-size="48" font-weight="bold">Итерация %d</text>
  `, path, color, iter)
}

func bounds(p []coastline.LatLon) (float64, float64, float64, float64) {
	minLat, maxLat := p[0].Lat, p[0].Lat
	minLon, maxLon := p[0].Lon, p[0].Lon
	for _, pt := range p {
		if pt.Lat < minLat { minLat = pt.Lat }
		if pt.Lat > maxLat { maxLat = pt.Lat }
		if pt.Lon < minLon { minLon = pt.Lon }
		if pt.Lon > maxLon { maxLon = pt.Lon }
	}
	return minLat, maxLat, minLon, maxLon
}