// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"image"
	"image/color"
	"log"
	"terminal/internal/framebuffer"
	"terminal/internal/touchscreen"
	"time"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

func main() {
	// Инициализируем фреймбуффер
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		log.Panic(err)
	}

	// Получаем смапленное в память фреймбуффера изображение
	display, err := fb.Image()
	if err != nil {
		log.Panic(err)
	}

	// Получаем устройство ввода
	touch, err := touchscreen.New("/dev/input/event0")
	if err != nil {
		log.Panic(err)
	}
	defer touch.Close()

	// Создаем GUI
	gui, err := sgui.New(display, &touch)
	if err != nil {
		panic(err)
	}

	// Создаем тему
	theme := widget.ColorTheme{
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		MainColor:       color.RGBA{200, 200, 200, 255},
		SecondColor:     color.RGBA{180, 180, 180, 255},
		StrokeColor:     color.RGBA{60, 60, 60, 255},
		TextColor:       color.Black,
		StrokeWidth:     2,
		CornerRadius:    10,
	}

	// Создаем экраны
	mainScreen := sgui.NewScreen(gui.SizeDisplay())
	mainScreen.SetBackground(theme.BackgroundColor)

	secondScreen := sgui.NewScreen(gui.SizeDisplay())
	secondScreen.SetBackground(theme.BackgroundColor)

	// Создаем виджеты на основной экран
	ind := widget.NewIndicator(20, theme)
	ind.AddState(color.RGBA{255, 0, 0, 255})
	ind.AddState(color.RGBA{0, 255, 0, 255})

	button2 := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				if ind.GetState() == 0 {
					ind.SetState(1)
				} else {
					ind.SetState(0)
				}
			},
			Label:           "Button 2",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	button1 := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				if button2.Hidden() {
					button2.Show()
				} else {
					button2.Hide()
				}
			},
			Label:           "Hide",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	buttonSetSecondScreen := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				gui.SetScreen(&secondScreen)
			},
			Label:           "2 экран",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	buttonSetMainScreen := widget.NewButton(
		widget.ButtonParam{
			Size: image.Point{X: 110, Y: 40},
			Onclick: func() {
				gui.SetScreen(&mainScreen)
			},
			Label:           "1 экран",
			LabelSize:       20,
			ReleaseColor:    theme.MainColor,
			PressColor:      theme.SecondColor,
			BackgroundColor: theme.BackgroundColor,
			CornerRadius:    theme.CornerRadius,
			StrokeWidth:     theme.StrokeWidth,
			StrokeColor:     theme.StrokeColor,
			TextColor:       theme.TextColor,
		},
	)

	// Добавляем виджеты на холст
	mainScreen.AddWidget(10, 10, button1)
	mainScreen.AddWidget(10, 60, button2)
	mainScreen.AddWidget(130, 70, ind)
	mainScreen.AddWidget(10, 200, buttonSetSecondScreen)

	secondScreen.AddWidget(10, 10, buttonSetMainScreen)

	// Устанавливаем активный экран
	gui.SetScreen(&mainScreen)

	gui.StartInputEventHandler()

	for {

		//start := time.Now()
		gui.Render()
		//log.Printf("Rendering  %v\n", time.Since(start))

		time.Sleep(50 * time.Millisecond)
	}

}
