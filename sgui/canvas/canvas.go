package canvas

import "image"

// Позиция относительно базового объекта
type Position struct {
	X uint // Позиция от левой грани
	Y uint // Позиция верхней грани
}

type Size struct {
	Width  int
	Height int
}

type Widget interface {
	// Отрисовывает виджет
	Render() *image.RGBA

	GetSize() Size

	// Обработка нажатия и отпускания
	Tap()
	Release()
}

type CanvasObject interface {
	Widget
	SetPosition(Position)
}

type Canvas struct {
	Display *image.Image
	Objects []CanvasObject
}

func New(Display *image.Image) Canvas {
	return Canvas{Display: Display}
}
