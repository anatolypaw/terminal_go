package guiview

import (
	"image/color"
	"log"
	"terminal/internal/app"
	"terminal/internal/framebuffer"
	"terminal/internal/touchscreen"
	"time"

	"github.com/anatolypaw/sgui"
	"github.com/anatolypaw/sgui/widget"
)

type GuiView struct {
	sgui  *sgui.Sgui
	theme widget.ColorTheme

	// Экраны
	ScreenProduceCamera *screenProduceCamera
	screenSelectMode    *screenSelectMode
}

func New(a *app.App) GuiView {
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

	// Создаем GUI
	gui, err := sgui.New(display, &touch)
	if err != nil {
		panic(err)
	}

	// Создаем тему
	theme := widget.ColorTheme{
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		MainColor:       color.RGBA{230, 230, 230, 255},
		SecondColor:     color.RGBA{180, 180, 180, 255},
		StrokeColor:     color.RGBA{60, 60, 60, 255},
		TextColor:       color.Black,
		StrokeWidth:     1,
		CornerRadius:    6,
	}

	gv := GuiView{
		sgui:  &gui,
		theme: theme,
	}

	// Инциализируем экраны
	gv.initScreenMain()
	gv.initScreenSelectMode(a)

	return gv
}

func (v *GuiView) Run(a *app.App) {
	v.sgui.SetScreen(v.ScreenProduceCamera.Screen)

	v.sgui.StartInputEventHandler()

	// Обновляем данные в виджетах и делаем рендеринг
	for {
		v.ScreenProduceCamera.modeIndicator.SetState(a.I % 2)

		start := time.Now()
		v.sgui.Render()

		since := time.Since(start)

		if since > 10*time.Millisecond {
			log.Printf("Rendering  %v\n", since)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
