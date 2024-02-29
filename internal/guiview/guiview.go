package guiview

import (
	"image/color"
	"log"
	o2i500 "terminal/internal/O2i500"
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

	// Общие виджеты

	// Экраны
	ScreenProduceCamera *sgui.Screen
	ScreenSelectMode    *sgui.Screen
	ScreenSelecGood     *sgui.Screen
}

func New(app *app.App, o2i500 *o2i500.O2i500) GuiView {
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
	gv.ScreenProduceCamera = NewScreenMain(&gv, app, o2i500)
	gv.ScreenSelectMode = NewScreenSelectMode(&gv, app)
	gv.ScreenSelecGood = NewScreenSelectGood(&gv, app)

	return gv
}

func (v *GuiView) Run() {
	v.sgui.SetScreen(v.ScreenProduceCamera)
	v.sgui.StartInputEventHandler()

	// Обновляем данные в виджетах и делаем рендеринг
	for {
		start := time.Now()
		v.sgui.Render()

		since := time.Since(start)

		if since > 1*time.Millisecond {
			log.Printf("Rendering  %v\n", since)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
