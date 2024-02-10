package painter

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestDrawCircle(t *testing.T) {
	tests := []struct {
		name string
		c    Circle
		r    int
	}{
		{"red", Circle{FillColor: color.RGBA{255, 0, 0, 255}}, 45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fname := fmt.Sprintf("./test_result/%s.png", tt.name)
			f, err := os.Create(fname)
			if err != nil {
				return
			}
			defer f.Close()

			size := tt.r * 2
			r := image.Rect(0, 0, size, size)
			img := image.NewRGBA(r)

			DrawCircle(img, tt.c)

			err = png.Encode(f, img)
			if err != nil {
				return
			}

		})
	}
}
