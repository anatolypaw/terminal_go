package guiview

import (
	"fmt"
	"image"
	"image/color"
	"terminal/internal/app"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

// Это экран калибровки тачскрина
func NewScreenTouchCalib(gv *GuiView, a *app.App) *sgui.Screen {
	// Создаем экран
	s := sgui.NewScreen(gv.sgui.SizeDisplay())

	// Установка цвета фона
	s.SetBackground(color.Black)

	// Надпись "Калибровка экрана"
	textMode := widget.NewLabel(
		&widget.LabelParam{
			Size:            image.Point{300, 40},
			Text:            "Калибровка экрана",
			TextSize:        30,
			TextColor:       color.White,
			BackgroundColor: color.Black,
		},
		nil,
	)

	s.AddWidget(250, 10, textMode)

	// Перехватывает нажатие и позицию нажатия
	f := func(pos image.Point) {
		fmt.Println(pos)
	}

	taphooker := widget.NewTapHooker(f)

	s.AddWidget(0, 0, taphooker)

	return &s

}
