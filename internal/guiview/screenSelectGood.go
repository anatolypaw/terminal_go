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

	bparam := widget.ButtonParam{
		Size:             image.Point{300, 70},
		OnClick:          nil,
		Text:             "",
		TextSize:         30,
		ReleaseFillColor: gv.theme.MainColor,
		PressFillColor:   gv.theme.SecondColor,
		BackgroundColor:  gv.theme.BackgroundColor,
		CornerRadius:     gv.theme.CornerRadius,
		StrokeWidth:      gv.theme.StrokeWidth,
		StrokeColor:      gv.theme.StrokeColor,
		TextColor:        color.Black,
		Hidden:           false,
	}

	for i, g := range a.Goods {
		button := *widget.NewButton(
			nil,
			func() widget.ButtonParam {
				if g.Gtin == "" {
					bparam.Hidden = true
					bparam.Text = "---"
					return bparam
				}
				bparam.OnClick = func() {
					a.SelectedGood = g
					gv.sgui.SetScreen(gv.ScreenProduceCamera)
				}

				bparam.Hidden = false
				bparam.Text = g.Desc
				bparam.ReleaseFillColor = g.Color
				return bparam
			},
		)

		x := 60
		y := 85*((i%5)+1) - 25

		if i > 4 {
			x = 420
		}

		s.AddWidget(x, y, &button)
	}

	return &s
}
