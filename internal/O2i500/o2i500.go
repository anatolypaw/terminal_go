package o2i500

import (
	"fmt"
	"net"
	"os"
)

type O2i500 struct {
}

func New() O2i500 {

	return O2i500{}
}

func (o *O2i500) Run(address string) {
	// Установка соединения с сервером

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Буфер для принимаемых данных
	buffer := make([]byte, 1024)

	for {
		// Чтение данных из соединения
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Ошибка чтения данных:", err)
			return
		}

		// Вывод принятых данных
		fmt.Print(string(buffer[:n]))
	}
}

func (o *O2i500) ReadCode() {

}
