// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package framebuffer

import (
	"image"
	"image/color"
)

type BGRA struct {
	Pix    []byte
	Rect   image.Rectangle
	Stride int
}

func (i *BGRA) Bounds() image.Rectangle { return i.Rect }
func (i *BGRA) ColorModel() color.Model { return color.RGBAModel }

func (i *BGRA) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(i.Rect)) {
		return color.RGBA{}
	}

	n := i.PixOffset(x, y)
	pix := i.Pix[n:]
	return color.RGBA{pix[2], pix[1], pix[0], pix[3]}
}

func (i *BGRA) Set(x, y int, c color.Color) {
	r, g, b, a := c.RGBA()

	//	if !(image.Point{x, y}.In(i.Rect)) {
	//		return
	//	}

	n := i.PixOffset(x, y)
	pix := i.Pix[n : n+4]

	pix[0] = byte(b)
	pix[1] = byte(g)
	pix[2] = byte(r)
	pix[3] = byte(a)
}

func (i *BGRA) PixOffset(x, y int) int {
	return (y-i.Rect.Min.Y)*i.Stride + (x-i.Rect.Min.X)*4
}
