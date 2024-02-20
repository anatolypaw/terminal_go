package guiview

import (
	"image"
	"image/color"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

type screenProduceCamera struct {
	Screen *sgui.Screen

	// виджеты
	Counter       *widget.Label
	modeIndicator *widget.TextIndicator
}

func (gv *GuiView) initScreenMain() {
	// Создаем экран
	s := sgui.NewScreen(gv.sgui.SizeDisplay())

	// Установка цвета фона
	s.SetBackground(gv.theme.BackgroundColor)

	// Надпись "Режимм:"
	textMode := widget.NewLabel(
		widget.LabelParam{
			Size:            image.Point{65, 30},
			Text:            "Режим:",
			TextSize:        20,
			TextColor:       color.Black,
			BackgroundColor: gv.theme.BackgroundColor,
		})

	s.AddWidget(10, 10, textMode)

	// Индикатор режима
	modeIndicator := *widget.NewTextIndicator(
		widget.TextIndicatorParam{
			Size:            image.Point{300, 60},
			BackgroundColor: gv.theme.BackgroundColor,
			StrokeWidth:     gv.theme.StrokeWidth,
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
				gv.sgui.SetScreen(gv.screenSelectMode.Screen)
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

	gv.ScreenProduceCamera = &screenProduceCamera{
		Screen:        &s,
		modeIndicator: &modeIndicator,
	}
}

// Обновляет данные на экране
func (s *screenProduceCamera) UpdateData() {

}
