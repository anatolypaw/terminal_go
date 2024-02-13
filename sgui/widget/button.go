package widget

import (
	"image"
	"image/color"
	"image/draw"
	"sgui/entity"
	"sgui/widget/painter"
	"sgui/widget/painter/text2img"
)

// Имеет неограниченное количество переключающихся состояний
// Состояние представляет собой круг заданного цвета
// 1) Сначала нужно создать состояние через AddState()
// 2) Для изменения состояния испольузется SetState()

type button struct {
	size         entity.Size
	onClick      func()
	currentState int  // Текущее состояние
	tapped       bool // Флаг, что кнопка нажата

	releasedRender *image.RGBA
	pressedRenser  *image.RGBA
}

type Button struct {
	Size      entity.Size
	Onclick   func()
	Label     string
	LabelSize float64
}

func NewButton(b Button) *button {
	if b.Size.Height <= 0 {
		b.Size.Height = 1
	}

	if b.Size.Width <= 0 {
		b.Size.Width = 1
	}

	// Состояния кнопки.
	// 0 - кнопка отжата
	// 1 - кнопка нажата
	releasedRender := painter.DrawRectangle(painter.Rectangle{
		Width:        b.Size.Width,
		Height:       b.Size.Height,
		FillColor:    color.RGBA{94, 94, 94, 255},
		CornerRadius: 8,
		StrokeWidth:  1,
		StrokeColor:  color.RGBA{34, 34, 34, 255},
	})
	pressedRender := painter.DrawRectangle(painter.Rectangle{
		Width:        b.Size.Width,
		Height:       b.Size.Height,
		FillColor:    color.RGBA{118, 118, 118, 255},
		CornerRadius: 8,
		StrokeWidth:  1,
		StrokeColor:  color.RGBA{34, 34, 34, 255},
	})

	// Получаем изображение текста и вычисляем его расположение
	// для размещения в середине кнопки
	textimg := text2img.Text2img(b.Label, b.LabelSize)
	textMidPos := image.Point{
		X: -(b.Size.Width - textimg.Rect.Dx()) / 2,
		Y: -(b.Size.Height - textimg.Rect.Dy()) / 2,
	}

	// Наносим текст на оба состояния кнопки
	draw.Draw(releasedRender,
		releasedRender.Bounds(),
		textimg,
		textMidPos,
		draw.Over)

	draw.Draw(pressedRender,
		pressedRender.Bounds(),
		textimg,
		textMidPos,
		draw.Over)

	return &button{
		size: entity.Size{
			Width:  b.Size.Width,
			Height: b.Size.Height,
		},
		onClick:        b.Onclick,
		releasedRender: releasedRender,
		pressedRenser:  pressedRender,
	}
}

func (w *button) Render() *image.RGBA {
	if w.tapped {
		return w.pressedRenser
	} else {
		return w.releasedRender
	}
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
