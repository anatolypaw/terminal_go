package touchscreen

import (
	"bytes"
	"encoding/binary"
	"image"
	"log"
	"os"
	"time"

	"github.com/anatolypaw/sgui"
)

const (
	TYPE_SYNC  = 0 //
	TYPE_PRESS = 1 // Нажатие на тач.
	TYPE_ABS   = 3 // Координаты нажатия

	CODE_FORCE = 24 // усилие нажатия
	CODE_X     = 0  // х координата
	CODE_Y     = 1  // y координата

)

type Touchscreen struct {
	device  *os.File
	inputb  [24]byte
	pressed bool // флаг, состояние нажатия
	ptrig   bool // триггер нажатия
}

func New(path string) (Touchscreen, error) {
	f, err := os.Open(path)
	if err != nil {
		return Touchscreen{}, err
	}

	return Touchscreen{device: f}, nil
}

func (t *Touchscreen) Close() {
	t.device.Close()
}

// Возвращает событие нажатия или отпускания кнопки
// Событие смены координаты не реализовано
func (t *Touchscreen) GetEvent() sgui.IEvent {

	var event uint8 // Тип события
	var code uint8  // Код события
	var value int16 // Значение события

	var x int
	var y int
	var max_force int // Записывается максимальное зарегистрированное усилие

	// Увеличивается при поступлении соответсвующего события
	var xcount int // счетчик событий координат X
	var ycount int // счетчик событий координат Y

	for {
		// Читаем события в буффер, пока не получим все три события:
		// Нажатие, координата Х, координата Y
		// В одном чтении приходит одно событие
		_, err := t.device.Read(t.inputb[:])
		if err != nil {
			log.Println("ERR Touchscreen: ", err)
			time.Sleep(1 * time.Second)
		}

		//fmt.Printf("%3v\n", t.inputb)

		// Время события
		/*
			sec := binary.LittleEndian.Uint16(inputb[0:4])
			usec := binary.LittleEndian.Uint16(inputb[4:8])
			time := time.Unix(int64(sec), int64(usec))
		*/

		event = uint8(binary.LittleEndian.Uint16(t.inputb[8:10]))
		code = uint8(binary.LittleEndian.Uint16(t.inputb[10:12]))

		binary.Read(bytes.NewReader(t.inputb[12:14]), binary.LittleEndian, &value)

		// Событие координаты х
		if event == TYPE_ABS && code == CODE_X {
			x = int(value)
			xcount++
		}

		// Событие координаты y
		if event == TYPE_ABS && code == CODE_Y {
			y = int(value)
			ycount++
		}

		// событие усилия нажатия
		// сохраняем максимальное усилие
		if event == TYPE_ABS && code == CODE_FORCE {
			if uint8(value) > uint8(max_force) {
				max_force = int(value)
			}
		}

		// Событие нажатия
		// Нажато: 1, отпущено: 0
		if event == TYPE_PRESS && value == 1 {
			t.ptrig = true
		}

		// Событие отпускания.
		if event == TYPE_PRESS && value == 0 && t.pressed {
			t.pressed = false
			return sgui.EventRelease{
				Pos: image.Point{X: x, Y: y},
			}
		}

		// Возвращаем нажатие
		// Сработка нажатия только при сильном нажатии на тач
		if t.ptrig && max_force > 150 && xcount > 2 && ycount > 2 {
			t.ptrig = false
			t.pressed = true
			return sgui.EventTap{
				Pos: image.Point{X: x, Y: y},
			}
		}

	}

}
