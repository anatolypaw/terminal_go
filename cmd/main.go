// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"image/color"
	"log"
	"sgui"
	"sgui/widget"
	"terminal/framebuffer"
	"terminal/touchscreen"
	"time"
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
		log.Panic(err)
	}
	gui.StartInputWorker()

	button := widget.NewButton(50, 50, "click me", nil)
	button2 := widget.NewButton(200, 300, "ousshhit", nil)

	ind := widget.NewIndicator(30)
	ind.AddState(color.RGBA{0, 0, 255, 255})
	ind.AddState(color.RGBA{255, 0, 0, 255})
	ind.AddState(color.RGBA{255, 255, 0, 255})
	ind.AddState(color.RGBA{255, 0, 0, 255})

	gui.AddWidget(400, 40, button)
	gui.AddWidget(400, 100, button2)

	gui.AddWidget(465, 50, ind)

	var i int
	for {
		time.Sleep(1 * time.Second)
		ind.SetState(i % ind.States())
		gui.Render()
		i++
	}

}
