package painter

import (
	"image/color"
	"image/draw"

	"github.com/srwiley/rasterx"
)

// Circle describes a colored circle primitive in a Fyne canvas
type Circle struct {
	FillColor color.Color // The circle fill color
}

// DrawFastCircle rasterizes the given circle object into an image.
func DrawFastCircle(img draw.Image, c Circle) {

	r := img.Bounds().Dx() / 2
	x, y, dx, dy := r-1, 0, 1, 1

	err := dx - (r * 2)
	x0 := r
	y0 := x0

	for x > y {
		img.Set(x0+x, y0+y, c.FillColor)
		img.Set(x0+y, y0+x, c.FillColor)
		img.Set(x0-y, y0+x, c.FillColor)
		img.Set(x0-x, y0+y, c.FillColor)
		img.Set(x0-x, y0-y, c.FillColor)
		img.Set(x0-y, y0-x, c.FillColor)
		img.Set(x0+y, y0-x, c.FillColor)
		img.Set(x0+x, y0-y, c.FillColor)

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

func DrawCircle(img draw.Image, c Circle) {

	radius := img.Bounds().Dx() / 2

	scanner := rasterx.NewScannerGV(img.Bounds().Dx(), img.Bounds().Dy(), img, img.Bounds())

	filler := rasterx.NewFiller(img.Bounds().Dx(), img.Bounds().Dy(), scanner)
	filler.SetColor(c.FillColor)

	rad := float64(radius)
	rasterx.AddCircle(rad, rad, rad*0.95, filler)
	filler.Draw()

}
