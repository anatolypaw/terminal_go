package savema

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Savema struct {
	Address string
	conn    net.Conn
	Online  bool

	sendmu          sync.Mutex      // Для блокировки одновременного вызова Send
	lastSendtime    time.Time       // Время, когда была последня отправка сообщения в принтер
	expectedCommand expectedCommand // Тип команды, который ожидаем получить в ответе

	respValue chan string // в этот канад пересылается ответ на отправленную команду
	ioerror   chan error  //

}

// Команда, которую ожидаем считать. Потокобезопасная
type expectedCommand struct {
	command string
	mu      sync.Mutex
}

func (e *expectedCommand) Set(command string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.command = command
}

func (e *expectedCommand) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.command = ""
}

func (e *expectedCommand) Get() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.command
}

func New(address string) Savema {
	return Savema{
		Address:   address,
		ioerror:   make(chan error),
		respValue: make(chan string),
	}
}

func (s *Savema) Run() {
	for {
		s.Online = false
		log.Print("Установка соединения с принтером ", s.Address)
		err := s.connect()
		if err != nil {
			// Если не удается установить соединение, прервать цикл и повторить попытку.
			log.Print("Ошибка подключения к принтеру:", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		s.Online = true
		log.Print("Соединение с принтером установлено")

		// Запускаем чтение сообщений от принтера
		go s.read()

		fmt.Println(s.Send("SPGGSN", "")) // Вывести серийный номер

		fmt.Println(s.Send("SPCSPM", "1>PRINTED")) // Включить отчет о печати

		// Ждем, пока не случится ошибка обмена данными с савемой.
		// Если произойдет, перезапускаем соединение
		<-s.ioerror
		s.conn.Close()
	}
}

// Установка соединения с сервером
func (s *Savema) connect() error {
	c, err := net.DialTimeout("tcp", s.Address, 5*time.Second)
	if err != nil {
		return err
	}
	s.conn = c
	return nil
}

// Получаем сообещния от принтера
func (s *Savema) read() {
	// буффер для пакета
	var buffer [4096]byte

	// Паттерн ответа принтера
	// Ожидается сообщение в формате ~SPGRES{SPGGSN:22050235}^
	regxpResp := regexp.MustCompile(`~SPGRES\{([^{}]*)\}\^`)

	// Паттерн для разбивки ответа принтера на команду и значение
	// ожидается сообщение в формате SPGGSN:22050235
	// regxpCommand := regexp.MustCompile(`^(\w*):(.*)`)

	for {
		b, err := s.conn.Read(buffer[:])
		if err != nil {
			log.Print("Savema: ошибка чтения данных:", err)
			s.ioerror <- err
			return
		}
		str := string(buffer[:b])
		fmt.Println("<< RAW: ", str)

		// Поиск всех сообщений принтера
		matches := regxpResp.FindAllStringSubmatch(str, -1)

		for _, matches := range matches {
			for i, match := range matches {
				// Первый элемент - полное совпадение. пропускаем его.
				if i == 0 {
					continue
				}

				s.processResponse(match)
			}
		}
	}
}

// Отправляет команду принтеру, возвращает ответ и сколько принтер обрабатывал команду
func (s *Savema) Send(command string, value string) (string, time.Duration, error) {
	s.sendmu.Lock()
	defer s.sendmu.Unlock()
	var msg string
	if value == "" {
		msg = "~" + command + "^"
	} else {
		msg = "~" + command + "{" + value + "}^"
	}

	// Делаем задержку между отправками, иначе принтер не примет сообщение
	const delay = 10 * time.Millisecond    // Задержка между отправкой пакетов
	const timeout = 100 * time.Millisecond // Таймаут получения ответа на запрос

	nextTime := s.lastSendtime.Add(delay)
	time.Sleep(time.Until(nextTime))

	s.lastSendtime = time.Now()
	_, err := s.conn.Write([]byte(msg))
	if err != nil {
		s.ioerror <- err
		return "", time.Since(s.lastSendtime), err
	}

	// Фиксируем, от какой команды ожидаем ответ
	s.expectedCommand.Set(command)
	defer s.expectedCommand.Reset()

	// Ждем ответа принтера. ответ должен содержать имя команды.
	timeIsUp := make(chan bool)

	go func() {
		time.Sleep(timeout)
		timeIsUp <- true
	}()

	// Ждем получение этой команды за заданное время
	var respValue string
	select {
	case <-timeIsUp:
		return "", time.Since(s.lastSendtime), errors.New("превышено время ожидания ответа от принтера")

	case respValue = <-s.respValue:
		return respValue, time.Since(s.lastSendtime), nil
	}

}

// Обрабатывает полученный от принтера ответ
func (s *Savema) processResponse(resp string) {
	// Ожидаемый ответ состоит из двух частей - команды и значения, разделенных :
	// Но может прийти сообещние без разделителя
	// Разделим его на эти части
	// Разделение строки на две части по первому вхождению разделителя
	parts := strings.SplitN(resp, ":", 2)
	var command string
	var value string
	if len(parts) == 2 {
		command = parts[0]
		value = parts[1]
	} else {
		command = resp
	}

	switch command {
	// Ожидаемый ответ на запрос
	case s.expectedCommand.Get():
		// неблокирующая отправка в канал, который читает функция отправки команды
		select {
		case s.respValue <- value:
		default:
		}

	// Асинхронные сообщения
	// Включен отчет принтера о печати
	case "SPCSPM:OK":
		fmt.Println("Включен отчет о печати")

	default:
		fmt.Println("Неизвестный ответ принтера:", command, value)
	}

}
