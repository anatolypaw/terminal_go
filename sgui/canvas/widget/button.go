package widget

import (
	"canvas/widget/painter"
	"image"
	"image/color"
)

// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type buttonState struct {
	img *image.RGBA
}

type button struct {
	height int
	width  int
	label  string

	onClick func()

	currentState int  // Текущее состояние
	tapped       bool // Флаг, что кнопка нажата
	states       []buttonState
}

func NewButton(height int, width int, label string, onClick func()) button {
	if height <= 0 {
		height = 1
	}

	if width <= 0 {
		width = 1
	}

	// Состояния кнопки.
	// 0 - кнопка отжата
	// 1 - кнопка нажата
	states := []buttonState{
		{img: painter.DrawRectangle(painter.Rectangle{
			Width:        width,
			Height:       height,
			FillColor:    color.RGBA{94, 94, 94, 255},
			CornerRadius: 5,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		})},
		{img: painter.DrawRectangle(painter.Rectangle{
			Width:        width,
			Height:       height,
			FillColor:    color.RGBA{118, 118, 118, 255},
			CornerRadius: 5,
			StrokeWidth:  1,
			StrokeColor:  color.RGBA{34, 34, 34, 255},
		})},
	}

	return button{
		height:  height,
		width:   width,
		label:   label,
		onClick: onClick,
		states:  states,
	}
}

func (w *button) Render() *image.RGBA {
	return w.states[w.currentState].img
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
