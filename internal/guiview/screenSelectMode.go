package guiview

import (
	"image"
	"image/color"
	"terminal/internal/app"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

// Это экран выбора режимов работы.
// Режим работы
// 0 - Производство
// 1 - Отбраковка
// 2 - Смена даты
// 3 - Информация

type screenSelectMode struct {
	Screen *sgui.Screen

	// Виджеты
}

func (gv *GuiView) initScreenSelectMode(a *app.App) {
	// Создаем экран
	s := sgui.NewScreen(gv.sgui.SizeDisplay())

	// Установка цвета фона
	s.SetBackground(gv.theme.BackgroundColor)

	// Надпись "Выбор режима"
	textMode := widget.NewLabel(
		widget.LabelParam{
			Size:            image.Point{230, 40},
			Text:            "Выбор режима",
			TextSize:        30,
			TextColor:       color.Black,
			BackgroundColor: gv.theme.BackgroundColor,
		})

	s.AddWidget(280, 10, textMode)

	// Кнопка выбора режима производство
	buttonProduceMode := *widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{300, 60},
			Onclick: func() {
				a.SetMode(0)
				gv.sgui.SetScreen(gv.ScreenProduceCamera.Screen)
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
		})

	s.AddWidget(240, 100, &buttonProduceMode)

	// Кнопка выбора режима отбраковки
	buttonCancelMode := *widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{300, 60},
			Onclick: func() {
				a.SetMode(1)
				gv.sgui.SetScreen(gv.ScreenProduceCamera.Screen)
			},
			Label:           "ОТБРАКОВКА",
			LabelSize:       30,
			ReleaseColor:    gv.theme.MainColor,
			PressColor:      gv.theme.SecondColor,
			BackgroundColor: gv.theme.BackgroundColor,
			CornerRadius:    gv.theme.CornerRadius,
			StrokeWidth:     gv.theme.StrokeWidth,
			StrokeColor:     gv.theme.StrokeColor,
			TextColor:       color.Black,
		})
	s.AddWidget(240, 200, &buttonCancelMode)

	gv.screenSelectMode = &screenSelectMode{
		Screen: &s,
	}

}
