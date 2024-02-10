package widget

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
	// Возвращает текущую позицию виджета
	Position() Position

	// Отрисовывает виджет
	Rendewr() *image.RGBA
}
