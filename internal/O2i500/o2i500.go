package o2i500

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

type O2i500 struct {
	LastPacketTime time.Time
	Connected      bool
}

func New() O2i500 {

	return O2i500{}
}

func (o *O2i500) Run(address string) {
	for {

		// Установка соединения с сервером
		log.Print("Установка соединения с камерой ", address)
		conn, err := net.DialTimeout("tcp", address, 1000*time.Millisecond)
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}
		defer conn.Close()

		// Постоянная отправка данных, для проверки состояния соединения
		o.LastPacketTime = time.Now()

		go func() {
			for {
				time.Sleep(1000 * time.Millisecond)
				_, err := conn.Write([]byte("?"))
				if err != nil {
					conn.Close()
					return
				}

				t := o.LastPacketTime.Add(1000 * time.Millisecond)

				if t.Unix() < time.Now().Unix() {
					log.Print("Соединение с камерой умерло")
					conn.Close()
					o.Connected = false
					return
				}
				o.Connected = true
			}

		}()

		// буффер для пакета
		var buffer [4096]byte

		// Паттерн пакета с кодов
		// Две группы.
		// 0 - полное совпадение
		// 1 - Распознан ли код (0 или 1)
		// 2 - Считанный код
		r := regexp.MustCompile(`[0-9]{4}start;(.);(.*);stop`)

		for {
			b, err := conn.Read(buffer[:])
			if err != nil {
				fmt.Println("Ошибка чтения данных:", err)
				break
			}
			s := string(buffer[:b])

			//fmt.Println(s)
			o.LastPacketTime = time.Now()

			// Поиск всех вхождений с помощью регулярного выражения
			matches := r.FindAllStringSubmatch(s, -1)

			// Обработка найденных совпадений
			for _, match := range matches {

				// Код не распознан
				if match[1] == "0" {
					fmt.Println("Код не распознан")
					continue
				}

				fmt.Println("Camera: ", match[2])
			}

		}

		conn.Close()
	}
}

func (o *O2i500) ReadCode() {

}
