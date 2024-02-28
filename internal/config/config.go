package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type TouchCalib struct {
	Done    bool
	Point_1 int
	Point_2 int
	Point_3 int
	Point_4 int
}

type Config struct {
	TouchCalib   TouchCalib // Калибровки тачскрина
	TerminalType string     // Тип терминала. Камера или нанесение
	TerminalName string     // Имя терминала
	O2i500Addr   string     // Адрес камеры O2I500
	SavemaAddr   string     // Адрес савемы
	HubAddr      string     // Адрес и порт хаба
}

// Значения по умолчанию для конфигурации
var DefaultConfig = Config{
	TouchCalib:   TouchCalib{},
	TerminalName: "TEST",
	TerminalType: "camera",
	O2i500Addr:   "10.0.4.11:50010",
	SavemaAddr:   "10.0.0.1",
	HubAddr:      "10.0.4.20:3000",
}

// Функция для загрузки конфигурации из файла
func Load(filename string) (Config, error) {
	var config Config

	// Попытка чтения файла конфигурации
	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	// Декодируем содержимое файла JSON в структуру конфигурации
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	// Проверка правильности полей
	if config.TerminalName == "" {
		return Config{}, fmt.Errorf("не указано имя терминала")
	}

	if config.TerminalType != "camera" {
		return Config{}, fmt.Errorf("недопустимый тип терминала")
	}

	return config, nil
}

// Функция для сохранения конфигурации в файл
func Save(filename string, config Config) error {
	// Кодируем конфигурацию в формат JSON
	encodedConfig, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	// Записываем конфигурацию в файл
	err = os.WriteFile(filename, encodedConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}
