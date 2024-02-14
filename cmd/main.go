// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
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
	// Создаем гуй
	gui, err := sgui.New(display, &touch)
	if err != nil {
		panic(err)
	}
	backcolor := color.RGBA{50, 50, 50, 255}
	gui.SetBackground(backcolor)

	for i := 0; i < 5; i++ {
		for n := 0; n < 10; n++ {
			// Создаем виджеты

			ind := widget.NewIndicator(20, backcolor)
			ind.AddState(color.RGBA{0, 0, 255, 255})
			ind.AddState(color.RGBA{0, 255, 0, 255})

			button := widget.NewButton(
				widget.ButtonParam{
					Size: image.Point{X: 110, Y: 40},
					Onclick: func() {
						if ind.GetState() == 0 {
							ind.SetState(1)
						} else {
							ind.SetState(0)
						}
					},
					Label:     fmt.Sprintf("Button %v", n+(i*10)),
					LabelSize: 20,
				},
				backcolor)

			// Добавляем виджеты
			gui.AddWidget(10+i*160, 10+(n*47), button)
			gui.AddWidget(130+i*160, 20+(n*47), ind)

		}
	}

	gui.StartInputEventHandler()

	for {

		//start := time.Now()
		gui.Render()
		//log.Printf("Rendering  %v\n", time.Since(start))

		time.Sleep(100 * time.Millisecond)
	}

}
