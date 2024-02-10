// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"canvas/widget"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
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
	/*
		inputb := make([]byte, 24)

		var x int
		var y int
		var x_done bool
		var y_done bool
	*/
	indicator := widget.NewIndicator(10)
	ind := indicator.Render()
	/*	for {
		//f.Read(inputb)
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
				x = (event.value / -5) + 100
				x_done = true
			}
			if event.code == CODE_Y {
				y = event.value/8 - 450
				y_done = true
			}

		}

		if x_done && y_done && event.typee == TYPE_ABS {
			t := time.Now()

			draw.Draw(ui.Display, ui.Display.Bounds(), ind, image.Point{x, y}, draw.Src)
			fmt.Printf("Elapsed %v\n", time.Since(t))
			x_done = false
			y_done = false
		}

	} */
	for {

		img := draw.Image(ind)
		for i := 0; i < 800; i += 1 {
			draw.Draw(ui.Display, ui.Display.Bounds(), img, image.Point{i - 800, int(math.Sin(float64(i)/50)*20) - 200}, draw.Over)
		}

		for i := 800; i > 0; i -= 1 {
			draw.Draw(ui.Display, ui.Display.Bounds(), img, image.Point{i - 800, int(math.Sin(float64(i)/20)*100) - 200}, draw.Over)
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
