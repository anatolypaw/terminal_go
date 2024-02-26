package guiview

import (
	"image"
	"image/color"
	"terminal/internal/app"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

func NewScreenSelectGood(gv *GuiView, a *app.App) *sgui.Screen {
	// Создаем экран
	s := sgui.NewScreen(gv.sgui.SizeDisplay())

	// Установка цвета фона
	s.SetBackground(gv.theme.BackgroundColor)

	// Надпись "Выбор продукта:"
	textMode := widget.NewLabel(
		&widget.LabelParam{
			Size:            image.Point{230, 40},
			Text:            "Выбор продукта",
			TextSize:        30,
			TextColor:       color.Black,
			BackgroundColor: gv.theme.BackgroundColor,
		},
		nil,
	)

	s.AddWidget(280, 10, textMode)

	for _, g := range a.Goods {
		button := *widget.NewButton(
			widget.ButtonParam{
				Size: image.Point{300, 60},
				Onclick: func() {
					a.SetMode(app.MODE_PRODUCE)
					gv.sgui.SetScreen(gv.ScreenProduceCamera)
				},
				Label:           "ПРОИЗВОДСТВО",
				LabelSize:       30,
				ReleaseColor:    color.RGBA{200, 255, 200, 255},
				PressColor:      color.RGBA{110, 178, 140, 255},
				BackgroundColor: gv.theme.BackgroundColor,
				CornerRadius:    gv.theme.CornerRadius,
				StrokeWidth:     gv.theme.StrokeWidth,
				StrokeColor:     gv.theme.StrokeColor,
				TextColor:       color.Black,
			},
		)

		s.AddWidget(240, 100, &button)
		_ = g
		//
	}

	return &s
}
