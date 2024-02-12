package widget

import (
	"image"
	"image/color"
	"sgui/entity"
	"sgui/widget/painter"
)

// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type buttonRender struct {
	img *image.RGBA
}

type button struct {
	size  entity.Size
	label string

	onClick func()

	currentState int  // Текущее состояние
	tapped       bool // Флаг, что кнопка нажата
	renders      []buttonRender
}

func NewButton(width int, height int, label string, onClick func()) *button {
	if height <= 0 {
		height = 1
	}

	if width <= 0 {
		width = 1
	}

	// Состояния кнопки.
	// 0 - кнопка отжата
	// 1 - кнопка нажата
	renders := []buttonRender{
		{img: painter.DrawRectangle(painter.Rectangle{
			Width:        width,
			Height:       height,
			FillColor:    color.RGBA{94, 94, 94, 255},
			CornerRadius: 8,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		})},
		{img: painter.DrawRectangle(painter.Rectangle{
			Width:        width,
			Height:       height,
			FillColor:    color.RGBA{118, 118, 118, 255},
			CornerRadius: 8,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		})},
	}

	painter.AddLabel(renders[0].img, width/4, height/2, label)
	painter.AddLabel(renders[1].img, width/4, height/2, label)

	return &button{
		size: entity.Size{
			Width:  width,
			Height: height,
		},
		label:   label,
		onClick: onClick,
		renders: renders,
	}
}

func (w *button) Render() *image.RGBA {
	return w.renders[w.currentState].img
}

// Вызвать при нажатии на кнопку
func (w *button) Tap() {
	w.currentState = 1
	w.tapped = true
}

// Вызвать при отпускании кнопки
func (w *button) Release() {
	if w.tapped {
		w.Click()
	}
	w.currentState = 0
	w.tapped = false
}

// Вызвывается когда предварительно нажатая кнопка была отпущенна
func (w *button) Click() {

	if w.onClick != nil {
		go w.onClick()
	}
}

func (w *button) Size() entity.Size {
	return w.size
}
