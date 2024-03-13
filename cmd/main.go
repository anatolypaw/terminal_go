// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	"terminal/internal/app"
	"terminal/internal/config"
	"terminal/internal/guiview"
)

const version = "0.0.1"

func main() {
	exit := make(chan os.Signal, 1)

	configFileName := "config.json"

	// Парсим флаги командной строки
	newConfigFlag := flag.Bool("default", false, "создать config.json конфигурации по умолчанию.")
	flag.Parse()

	// Если указан параметр, создаем файл конфигурации по умолчанию
	if *newConfigFlag {
		err := config.Save(configFileName, config.DefaultConfig)
		if err != nil {
			log.Print("Ошибка при создании файла конфигурации:", err)
			return
		}
		log.Print("Создан файл конфигурации по умолчанию ", configFileName)
		return
	}

	log.Print("version ", version)

	// Загрузка конфигурации из файла
	cfg, err := config.Load(configFileName)
	if err != nil {
		log.Print("Ошибка загрузки конфигурации: ", err)
		return
	}
	log.Printf("Конфигурация загружена, тип: %s, имя: %s", cfg.TerminalType, cfg.TerminalName)

	app := app.New(cfg)
	go app.Run()

	// Запускаем графический интерфейс
	gui, err := guiview.New(&app)
	if err != nil {
		log.Print("Невозможно запустить GUI: ", err)
	} else {
		go gui.Run()
		log.Print("GUI запущен")
	}
	<-exit
}
