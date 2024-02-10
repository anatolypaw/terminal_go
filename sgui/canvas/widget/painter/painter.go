package painter

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

// Круг с растеризованными краями
func DrawCircle(size int, fillColor color.Color) *image.RGBA {
	rect := image.Rect(0, 0, size, size)
	img := image.NewRGBA(rect)

	radius := float64(size) / 2 * 0.93
	mid := float64(size) / 2

	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())

	// Рисуем основу
	filler := rasterx.NewFiller(size, size, scanner)
	filler.SetColor(fillColor)
	rasterx.AddCircle(mid, mid, radius, filler)
	filler.Draw()

	// Рисуем обводку
	dasher := rasterx.NewDasher(size, size, scanner)
	dasher.SetColor(color.RGBA{0, 0, 0, 100})
	dasher.SetStroke(fixed.Int26_6(size*3), 0, nil, nil, nil, 0, nil, 0)
	rasterx.AddCircle(mid, mid, radius, dasher)
	dasher.Draw()

	return img
}

func DrawFastCircle(img draw.Image, fillColor color.Color) {

	r := img.Bounds().Dx() / 2
	x, y, dx, dy := r-1, 0, 1, 1

	err := dx - (r * 2)
	x0 := r
	y0 := x0

	for x > y {
		img.Set(x0+x, y0+y, fillColor)
		img.Set(x0+y, y0+x, fillColor)
		img.Set(x0-y, y0+x, fillColor)
		img.Set(x0-x, y0+y, fillColor)
		img.Set(x0-x, y0-y, fillColor)
		img.Set(x0-y, y0-x, fillColor)
		img.Set(x0+y, y0-x, fillColor)
		img.Set(x0+x, y0-y, fillColor)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - r*2
		}
	}

}
