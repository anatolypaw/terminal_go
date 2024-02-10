package painter

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestDrawCircle(t *testing.T) {
	tests := []struct {
		name string
		c    color.Color
		size int
	}{
		{"red", color.RGBA{255, 0, 0, 255}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fname := fmt.Sprintf("./test_result/%s.png", tt.name)
			f, err := os.Create(fname)
			if err != nil {
				return
			}
			defer f.Close()

			img := DrawCircle(tt.size, tt.c)

			err = png.Encode(f, img)
			if err != nil {
				return
			}

		})
	}
}
