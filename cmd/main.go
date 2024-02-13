// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"image/color"
	"log"
	"sgui"
	"sgui/entity"
	"sgui/widget"
	"terminal/framebuffer"
	"terminal/touchscreen"
	"time"
)

func print(text string) {
	fmt.Println(text)
}

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
		log.Panic(err)
	}
	gui.StartInputWorker()

	ind := widget.NewIndicator(30)
	ind.AddState(color.RGBA{0, 0, 255, 255})
	ind.AddState(color.RGBA{0, 0, 0, 255})
	button := widget.NewButton(
		widget.Button{
			Size:      entity.Size{Width: 200, Height: 70},
			Onclick:   func() { log.Print("Событие кнопки 1") },
			Label:     "КН - 1",
			LabelSize: 50,
		},
	)

	var i int

	button2 := widget.NewButton(
		widget.Button{
			Size: entity.Size{Width: 60, Height: 30},
			Onclick: func() {
				log.Print("Событие кнопки 2")
				ind.SetState(i % ind.States())
				i++
				gui.Render()
			},
			Label:     "КН - 2",
			LabelSize: 10,
		},
	)

	gui.AddWidget(100, 100, button)
	gui.AddWidget(700, 400, button2)

	gui.AddWidget(465, 50, ind)

	for {
		time.Sleep(5 * time.Second)
	}

}
