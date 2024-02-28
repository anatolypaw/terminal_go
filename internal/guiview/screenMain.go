package guiview

import (
	"fmt"
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
		&widget.ButtonParam{
			Size: image.Point{60, 60},
			OnClick: func() {
				gv.sgui.SetScreen(gv.ScreenSelectMode)
			},
			Text:             "...",
			TextSize:         50,
			ReleaseFillColor: gv.theme.MainColor,
			PressFillColor:   gv.theme.SecondColor,
			BackgroundColor:  gv.theme.BackgroundColor,
			CornerRadius:     gv.theme.CornerRadius,
			StrokeWidth:      gv.theme.StrokeWidth,
			StrokeColor:      gv.theme.StrokeColor,
			TextColor:        color.Black,
		},
		nil,
	)

	s.AddWidget(315, 40, &modeMenu)

	// Индикатор выбранного продукта
	indicatorParam := widget.LabelParam{
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

	selectedProduct := *widget.NewLabel(nil,
		func() widget.LabelParam {

			good := app.SelectedGood

			if good.Gtin == "" {
				return indicatorParam
			}

			indicatorParam.Text = good.Desc
			indicatorParam.FillColor = good.Color

			return indicatorParam
		},
	)

	s.AddWidget(10, 110, &selectedProduct)

	// Кнопка выбора продукта
	// При нажатии открывает экран выбора продукта
	// Экран выбора продукта каждый раз создается новый
	// так как необходимо генерировать кнопки выбора продукта
	goodsMenu := *widget.NewButton(
		&widget.ButtonParam{
			Size: image.Point{60, 60},
			OnClick: func() {
				gv.sgui.SetScreen(gv.ScreenSelecGood)
			},
			Text:             "...",
			TextSize:         50,
			ReleaseFillColor: gv.theme.MainColor,
			PressFillColor:   gv.theme.SecondColor,
			BackgroundColor:  gv.theme.BackgroundColor,
			CornerRadius:     gv.theme.CornerRadius,
			StrokeWidth:      gv.theme.StrokeWidth,
			StrokeColor:      gv.theme.StrokeColor,
			TextColor:        color.Black,
		},
		nil,
	)

	s.AddWidget(315, 110, &goodsMenu)

	// Кнопка уменьшения даты
	dateDown := *widget.NewButton(
		nil,
		func() widget.ButtonParam {
			return widget.ButtonParam{
				Size: image.Point{60, 60},
				OnClick: func() {
					app.DateDown()
				},
				Text:             "-",
				TextSize:         50,
				ReleaseFillColor: gv.theme.MainColor,
				PressFillColor:   gv.theme.SecondColor,
				BackgroundColor:  gv.theme.BackgroundColor,
				CornerRadius:     gv.theme.CornerRadius,
				StrokeWidth:      gv.theme.StrokeWidth,
				StrokeColor:      gv.theme.StrokeColor,
				TextColor:        color.Black,
				Hidden:           !app.ModeIsProduce(),
			}
		},
	)

	s.AddWidget(10, 180, &dateDown)

	// Кнопка увеличения даты
	dateUp := *widget.NewButton(
		nil,
		func() widget.ButtonParam {
			return widget.ButtonParam{
				Size: image.Point{60, 60},
				OnClick: func() {
					app.DateUp()
				},
				Text:             "+",
				TextSize:         50,
				ReleaseFillColor: gv.theme.MainColor,
				PressFillColor:   gv.theme.SecondColor,
				BackgroundColor:  gv.theme.BackgroundColor,
				CornerRadius:     gv.theme.CornerRadius,
				StrokeWidth:      gv.theme.StrokeWidth,
				StrokeColor:      gv.theme.StrokeColor,
				TextColor:        color.Black,
				Hidden:           !app.ModeIsProduce(),
			}
		},
	)

	s.AddWidget(315, 180, &dateUp)

	//Дата производства
	labelProduceDate := *widget.NewLabel(nil,
		func() widget.LabelParam {
			return widget.LabelParam{
				Size:            image.Point{235, 60},
				Text:            fmt.Sprintf("23.10.20%2d", app.Date),
				TextSize:        40,
				TextColor:       color.Black,
				FillColor:       color.White,
				BackgroundColor: gv.theme.BackgroundColor,
				CornerRadius:    0,
				StrokeWidth:     gv.theme.StrokeWidth,
				StrokeColor:     gv.theme.StrokeColor,
				Hidden:          !app.ModeIsProduce(),
			}
		},
	)

	s.AddWidget(75, 180, &labelProduceDate)

	return &s
}
