// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"canvas/widget"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"os/signal"
	"sgui"
	"time"
)

const (
	TYPE_SYNC  = 0 //
	TYPE_PRESS = 1 // Нажатие на тач.
	TYPE_ABS   = 3 // Координаты нажатия

	CODE_FORCE = 24 // усилие нажатия
	CODE_X     = 0  // х координата
	CODE_Y     = 1  // y координата

)

type inputEvent struct {
	time  time.Time
	typee uint8
	code  uint8
	value int
}

func main() {
	ui, err := sgui.New()
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Open("/dev/input/event0")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	inputb := make([]byte, 24)

	var x int
	var y int
	var x_done bool
	var y_done bool

	ind := widget.NewIndicator(100)
	ind.AddState(color.RGBA{255, 0, 0, 255})
	ind.AddState(color.RGBA{0, 255, 0, 255})
	ind.AddState(color.RGBA{0, 0, 255, 255})

	var state int
	for {
		f.Read(inputb)
		sec := binary.LittleEndian.Uint16(inputb[0:4])
		usec := binary.LittleEndian.Uint16(inputb[4:8])

		var value int16
		binary.Read(bytes.NewReader(inputb[12:14]), binary.LittleEndian, &value)

		event := inputEvent{
			time:  time.Unix(int64(sec), int64(usec)),
			typee: uint8(binary.LittleEndian.Uint16(inputb[8:10])),
			code:  uint8(binary.LittleEndian.Uint16(inputb[10:12])),
			value: int(value),
		}

		if event.typee != TYPE_SYNC && event.typee == TYPE_ABS {

			if event.code == CODE_X {
				x = (event.value / -5) + 60
				x_done = true
			}
			if event.code == CODE_Y {
				y = event.value/8 - 480
				y_done = true
			}

		}

		if x_done && y_done && event.typee == TYPE_ABS && event.code == CODE_FORCE && event.value > 160 {
			t := time.Now()

			draw.Draw(ui.Display, ui.Display.Bounds(), ind.Render(), image.Point{x, y}, draw.Over)

			fmt.Printf("Elapsed %v\n", time.Since(t))
			x_done = false
			y_done = false
			ind.SetState(state % 6)
			state++
		}

	}
	// wait() // Wait until an exit signal has been received.
}

// wait polls for exit signals.
func wait() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	<-signals
}
