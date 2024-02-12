package sgui

import (
	"fmt"
	"image"
	"image/draw"
	"sgui/entity"
	"sgui/framebuffer"
	"time"
)

type sgui struct {
	fb      framebuffer.Framebuffer
	display draw.Image
	objects []Object
}

type IWidget interface {
	Render() *image.RGBA // Отрисовывает виджет
	Size() entity.Size
	Tap() // Обработка нажатия и отпускания
	Release()
}

type Object struct {
	Widget   IWidget
	Position image.Point
}

func New() (sgui, error) {
	// Инициализируем фреймбуффер
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		return sgui{}, err
	}

	// Получаем смапленное изображение
	display, err := fb.Image()
	if err != nil {
		return sgui{}, err
	}

	return sgui{
		fb:      *fb,
		display: display,
	}, nil
}

func (ui *sgui) Close() {
	ui.fb.Close()
}

// Возвращает размер дисплея
func (ui *sgui) DisplaySize() entity.Size {
	return entity.Size{
		Width:  ui.display.Bounds().Max.X,
		Height: ui.display.Bounds().Max.Y,
	}
}

// Добавляет объект (widget) на холст
func (ui *sgui) AddWidget(x int, y int, w IWidget) {
	obj := Object{
		Widget:   w,
		Position: image.Point{-x, -y},
	}
	ui.objects = append(ui.objects, obj)
}

// Отрисовывает объекты на дисплей
func (ui *sgui) Render() {

	start := time.Now()

	// Отрисовка на дисплей
	for _, o := range ui.objects {
		draw.Draw(
			ui.display,
			ui.display.Bounds(),
			o.Widget.Render(),
			o.Position,
			draw.Src)
	}

	fmt.Printf("Rendering  %v\n", time.Since(start))

}
