package widget

import (
	"canvas/widget/painter"
	"image"
	"image/color"
)

type BitIndicator struct {
	Radius int
	img    image.RGBA

	rendered bool // Флаг, что виджет был изменен и нужно его отрисовать

	// Логика виджета
	state bool
}

func NewIndicator(radius int) BitIndicator {
	size := radius * 2
	r := image.Rect(0, 0, size, size)
	w := BitIndicator{
		img: *image.NewRGBA(r),
	}
	return w
}

func (w *BitIndicator) Render() *image.RGBA {
	// Если вид индикатора не менялся, возвращаем текущее изображение
	if w.rendered {
		return &w.img
	}

	c := painter.Circle{
		// B G R
		FillColor: color.RGBA{255, 0, 0, 255},
	}

	painter.DrawCircle(&w.img, c)

	w.rendered = true
	return &w.img
}

// Логика виджета
func (w *BitIndicator) TurnOn() {
	if w.state {
		return
	}
	w.rendered = false
	w.state = true
}

// Включить
func (w *BitIndicator) TurnOff() {
	if !w.state {
		return
	}
	w.rendered = false
	w.state = false
}
