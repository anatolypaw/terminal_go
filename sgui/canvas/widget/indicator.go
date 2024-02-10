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

type bitIndicatorState struct {
	img *image.RGBA
}

type BitIndicator struct {
	size         int
	currentState int // Текущее состояние
	states       []bitIndicatorState
}

func NewIndicator(size int) BitIndicator {
	if size <= 0 {
		size = 1
	}
	return BitIndicator{size: size}
}

func (w *BitIndicator) AddState(c color.Color) {
	img := painter.DrawCircle(w.size, c)
	w.states = append(w.states, bitIndicatorState{img: img})
}

func (w *BitIndicator) SetState(s int) {
	if s < 0 {
		w.currentState = 0
		return
	}

	if s > len(w.states)-1 {
		w.currentState = len(w.states) - 1
		return
	}

	w.currentState = s

}

func (w *BitIndicator) Render() *image.RGBA {
	return w.states[w.currentState].img
}
