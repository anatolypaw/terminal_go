package painter

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/srwiley/rasterx"
)

func DrawCircle(size int, fillColor color.Color) *image.RGBA {
	rect := image.Rect(0, 0, size, size)
	img := image.NewRGBA(rect)

	radius := float64(size) / 2 * 0.9

	scanner := rasterx.NewScannerGV(img.Bounds().Dx(), img.Bounds().Dy(), img, img.Bounds())

	filler := rasterx.NewFiller(img.Bounds().Dx(), img.Bounds().Dy(), scanner)
	filler.SetColor(fillColor)

	rasterx.AddCircle(radius, radius, radius, filler)
	filler.Draw()

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
