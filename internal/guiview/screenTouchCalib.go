package guiview

import (
	"fmt"
	"image"
	"image/color"
	"terminal/internal/app"
	"time"

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
			Size:            image.Point{800, 40},
			Text:            "Калибровка тачскрина",
			TextSize:        30,
			TextColor:       color.White,
			BackgroundColor: color.Black,
		},
		nil,
	)
	s.AddWidget(0, 10, textMode)

	// Надпись "Калибровка экрана"
	label2 := widget.NewLabel(
		&widget.LabelParam{
			Size:            image.Point{800, 40},
			Text:            "Нажмите на экран для начала калибровки",
			TextSize:        20,
			TextColor:       color.White,
			BackgroundColor: color.Black,
		},
		nil,
	)
	s.AddWidget(0, 100, label2)

	// Перехватывает нажатие и позицию нажатия
	var step int // Шаг калибровки

	p1 := widget.NewRectangle(image.Point{10, 10}, color.White, color.Black)
	p1.Hide()

	p3 := widget.NewRectangle(image.Point{10, 10}, color.White, color.Black)
	p3.Hide()

	s.AddWidget(0, 0, p1)
	s.AddWidget(790, 470, p3)

	// Кнопка выхода
	button := widget.NewButton(&widget.ButtonParam{
		Size:             image.Point{20, 20},
		Text:             "ОК",
		TextSize:         10,
		ReleaseFillColor: gv.theme.MainColor,
		PressFillColor:   gv.theme.SecondColor,
		BackgroundColor:  color.Black,
		CornerRadius:     gv.theme.CornerRadius,
		TextColor:        color.Black,
		OnClick: func() {
			gv.sgui.SetScreen(gv.ScreenProduceCamera)
		},
	}, nil)

	button.Hide()
	s.AddWidget(390, 200, button)

	f := func(pos image.Point) {
		switch step {
		case 0:
			p1.Show()
			label2.SetText("Нажмите на белый квадрат в левом верхнем углу ", 20, color.White)
			time.Sleep(1 * time.Second)
		case 1:
			p1.Hide()
			p3.Show()
			label2.SetText("Нажмите на белый квадрат в правом нижнем углу", 20, color.White)
			time.Sleep(1 * time.Second)
		case 2:
			p3.Hide()
			label2.SetText("Готово", 20, color.White)
			button.Show()
		}

		step++
		fmt.Println(pos, step)
	}
	taphooker := widget.NewTapHooker(f)

	s.AddWidget(0, 0, taphooker)

	return &s

}
