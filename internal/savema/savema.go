package savema

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

type Savema struct {
	Address string
	conn    net.Conn
}

func New(address string) Savema {

	return Savema{
		Address: address,
	}
}

func (s *Savema) Run() {
	for {
		// Установка соединения с сервером
		log.Print("Установка соединения с принтером ", s.Address)
		c, err := net.DialTimeout("tcp", s.Address, 1000*time.Millisecond)
		if err != nil {
			log.Print("Ошибка подключения к принтеру:", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		s.conn = c

		log.Print("Соединение с принтером установлено")
		defer s.conn.Close()

		// буффер для пакета
		var buffer [4096]byte

		// Паттерн ответа принтера
		/* ([^{}]*) - это группа, которая соответствует любым символам,
		кроме { и } (потому что [^{}] соответствует любому символу,
		кроме { и }), и * указывает, что этот шаблон может повторяться любое
		 количество раз (включая ноль раз). 	*/
		regxp := regexp.MustCompile(`~SPGRES\{([^{}]*)\}\^`)

		s.Send("~SPGGSN^") // Вывести серийный номер
		time.Sleep(10 * time.Millisecond)
		s.Send("~SPCSPM{1>PRINTED}^") // Включить отчет о печати

		for {
			fmt.Print("<< ")
			b, err := s.conn.Read(buffer[:])
			if err != nil {
				log.Print("Savema: ошибка чтения данных:", err)
				break
			}
			str := string(buffer[:b])
			fmt.Println(str)

			// Поиск всех сообщений принтера
			matches := regxp.FindAllStringSubmatch(str, -1)

			for _, matches := range matches {
				for i, match := range matches {
					// Первый элемент - полное совпадение. пропускаем его.
					if i == 0 {
						continue
					}

					fmt.Println(match)
				}
			}

		}

		c.Close()
	}
}

// Отправляет команду принтеру, ждет ответа
func (s *Savema) Send(command string) (string, error) {
	fmt.Println(">>", command)
	_, err := s.conn.Write([]byte(command))
	return "", err
}
