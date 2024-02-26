package guiview

import (
	"image"
	"image/color"
	"terminal/internal/app"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

func NewScreenMain(gv *GuiView, app *app.App) *sgui.Screen {
	// Создаем экран
	s := sgui.NewScreen(gv.sgui.SizeDisplay())

	// Установка цвета фона
	s.SetBackground(gv.theme.BackgroundColor)

	// Надпись "Режимм:"
	textMode := widget.NewLabel(
		&widget.LabelParam{
			Size:            image.Point{65, 30},
			Text:            "Режим:",
			TextSize:        20,
			TextColor:       color.Black,
			BackgroundColor: gv.theme.BackgroundColor,
		},
		nil,
	)

	s.AddWidget(10, 10, textMode)

	// Индикатор режима производство / отбракова и тп
	modeIndicator := *widget.NewTextIndicator(
		widget.TextIndicatorParam{
			Size:            image.Point{300, 60},
			BackgroundColor: gv.theme.BackgroundColor,
			StrokeWidth:     gv.theme.StrokeWidth,
			StateSource:     app.GetMode,
		},
	)

	modeIndicator.AddState(
		"ПРОИЗВОДСТВО",
		30,
		color.Black,
		color.RGBA{200, 255, 200, 255},
		color.RGBA{110, 178, 140, 255},
	)

	modeIndicator.AddState(
		"ОТБРАКОВКА",
		30,
		color.Black,
		color.RGBA{208, 242, 253, 255},
		color.RGBA{101, 183, 209, 255},
	)

	s.AddWidget(10, 40, &modeIndicator)

	// Кнопка открытия экрана выбора режима.
	// При нажатии открывает экран выбора режимов
	modeMenu := *widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{60, 60},
			Onclick: func() {
				gv.sgui.SetScreen(gv.ScreenSelectMode)
			},
			Label:           "...",
			LabelSize:       50,
			ReleaseColor:    gv.theme.MainColor,
			PressColor:      gv.theme.SecondColor,
			BackgroundColor: gv.theme.BackgroundColor,
			CornerRadius:    gv.theme.CornerRadius,
			StrokeWidth:     gv.theme.StrokeWidth,
			StrokeColor:     gv.theme.StrokeColor,
			TextColor:       color.Black,
		})

	s.AddWidget(315, 40, &modeMenu)

	// Индикатор выбранного продукта
	selectedProduct := *widget.NewLabel(nil,
		func() widget.LabelParam {

			param := widget.LabelParam{
				Size:            image.Point{300, 60},
				Text:            "-----------",
				TextSize:        30,
				TextColor:       color.Black,
				FillColor:       color.White,
				BackgroundColor: gv.theme.BackgroundColor,
				CornerRadius:    0,
				StrokeWidth:     gv.theme.StrokeWidth,
				StrokeColor:     gv.theme.StrokeColor,
				Hidden:          false,
			}

			good := app.SelectedGood

			if good == nil {
				return param
			}

			param.Text = good.Desc
			param.FillColor = good.Color

			return param
		},
	)

	s.AddWidget(10, 110, &selectedProduct)

	// Кнопка выбора продукта
	// При нажатии открывает экран выбора продукта
	// Экран выбора продукта каждый раз создается новый
	// так как необходимо генерировать кнопки выбора продукта
	goodsMenu := *widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{60, 60},
			Onclick: func() {
				gv.sgui.SetScreen(gv.ScreenSelecGood)
			},
			Label:           "...",
			LabelSize:       50,
			ReleaseColor:    gv.theme.MainColor,
			PressColor:      gv.theme.SecondColor,
			BackgroundColor: gv.theme.BackgroundColor,
			CornerRadius:    gv.theme.CornerRadius,
			StrokeWidth:     gv.theme.StrokeWidth,
			StrokeColor:     gv.theme.StrokeColor,
			TextColor:       color.Black,
		})

	s.AddWidget(315, 110, &goodsMenu)

	return &s
}
