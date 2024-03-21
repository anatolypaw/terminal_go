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

	touch *touchscreen.Touchscreen
	// Общие виджеты

	// Экраны
	ScreenTouchCalib    *sgui.Screen
	ScreenProduceCamera *sgui.Screen
	ScreenSelectMode    *sgui.Screen
	ScreenSelecGood     *sgui.Screen
}

func New(app *app.App) (GuiView, error) {
	// Инициализируем фреймбуффер
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		return GuiView{}, err
	}

	// Получаем смапленное в память фреймбуффера изображение
	display, err := fb.Image()
	if err != nil {
		return GuiView{}, err
	}

	// Получаем устройство ввода
	touch, err := touchscreen.New("/dev/input/event0")
	if err != nil {
		return GuiView{}, err
	}

	// Создаем GUI
	gui, err := sgui.New(display, &touch)
	if err != nil {
		return GuiView{}, err
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
		touch: &touch,
	}

	// Инциализируем экраны
	gv.ScreenTouchCalib = NewScreenTouchCalib(&gv, app)
	gv.ScreenProduceCamera = NewScreenMain(&gv, app)
	gv.ScreenSelectMode = NewScreenSelectMode(&gv, app)
	gv.ScreenSelecGood = NewScreenSelectGood(&gv, app)

	return gv, nil
}

func (v *GuiView) Run() {
	// Сначала Переходим на экран калибровки тачскрина
	// Внутри этого экрана происходит проверка на наличие калибровки
	// Если она есть, то переключает на другой, зависящий от режима экран
	v.sgui.SetScreen(v.ScreenTouchCalib)
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
