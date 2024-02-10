package sgui

import (
	"image/draw"
	"sgui/framebuffer"
)

type sgui struct {
	Fb      framebuffer.Framebuffer
	Display draw.Image
}

func New() (sgui, error) {
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		return sgui{}, err
	}

	disp, err := fb.Image()
	if err != nil {
		return sgui{}, err
	}

	return sgui{Fb: *fb, Display: disp}, nil
}

func (ui *sgui) Close() {
	ui.Fb.Close()
}
