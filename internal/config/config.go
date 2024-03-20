package config

import (
	"encoding/json"
	"fmt"
	"os"
	"terminal/internal/entity"
)

type Config struct {
	filename string
	P        Params
}

type TouchCalib struct {
	Done    bool
	Point_1 int
	Point_2 int
	Point_3 int
	Point_4 int
}

type Params struct {
	TouchCalib   TouchCalib // Калибровки тачскрина
	TerminalName string     // Имя терминала
	TerminalType string     // Тип терминала. reading printing

	UseCamera  bool
	O2i500Addr string // Адрес камеры O2I500

	UseHandReader bool

	UseSavema  bool
	SavemaAddr string // Адрес савемы
	HubAddr    string // Адрес и порт хаба

	Goods [10]entity.Good
}

// Значения по умолчанию для конфигурации
var DefaultConfig = Params{
	TouchCalib:   TouchCalib{},
	TerminalName: "TEST",
	TerminalType: "printing",
	UseCamera:    false,
	O2i500Addr:   "10.0.4.11:50010",

	UseHandReader: true,
	UseSavema:     true,
	SavemaAddr:    "10.0.0.1",
	HubAddr:       "10.0.4.20:3100",
}

func New(file string) Config {
	return Config{
		filename: file,
		P:        Params{},
	}
}

// Функция для загрузки конфигурации из файла
func (c *Config) Load() error {
	var config Params

	// Попытка чтения файла конфигурации
	file, err := os.ReadFile(c.filename)
	if err != nil {
		return err
	}

	// Декодируем содержимое файла JSON в структуру конфигурации
	err = json.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	// Проверка правильности полей
	if config.TerminalName == "" {
		return fmt.Errorf("не указано имя терминала")
	}

	if config.TerminalType != "reading" && config.TerminalType != "printing" {
		return fmt.Errorf("недопустимый тип терминала")
	}

	c.P = config
	return nil
}

// Функция для сохранения конфигурации в файл
func (c *Config) Save() error {
	// Кодируем конфигурацию в формат JSON
	encodedConfig, err := json.MarshalIndent(c.P, "", "    ")
	if err != nil {
		return err
	}

	// Записываем конфигурацию в файл
	err = os.WriteFile(c.filename, encodedConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}
