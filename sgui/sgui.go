package sgui

import (
	"canvas"
	"sgui/framebuffer"
)

type sgui struct {
	Fb     framebuffer.Framebuffer
	Canvas canvas.Canvas
}

func New() (sgui, error) {
	// Инициализируем фреймбуффер
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		return sgui{}, err
	}

	img, err := fb.Image()
	if err != nil {
		return sgui{}, err
	}

	canvas := canvas.New(&img)

	return sgui{
		Fb:     *fb,
		Canvas: canvas,
	}, nil
}

func (ui *sgui) Close() {
	ui.Fb.Close()
}
