package touch

import "time"

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

/*
f, err := os.Open("/dev/input/event0")
	if err != nil {
		panic(err)
	}
	defer f.Close()


inputb := make([]byte, 24)
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

		if event.typee == TYPE_PRESS {
			t := time.Now()

			if event.value == 1 {
				button.Tap()
			} else {
				button.Release()
			}

			draw.Draw(ui.Display, ui.Display.Bounds(), button.Render(), image.Point{-400, -200}, draw.Src)
			fmt.Printf("Elapsed %v\n", time.Since(t))
		}
*/
