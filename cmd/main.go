// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"log"
	"os"
	o2i500 "terminal/internal/O2i500"
	"terminal/internal/app"
	"terminal/internal/config"
	"terminal/internal/guiview"
)

func main() {
	exit := make(chan os.Signal, 1)

	configFileName := "config.json"

	// Парсим флаги командной строки
	makeConfigFlag := flag.Bool(
		"make-config",
		false,
		"Будет создан файл config.json конфигурации по умолчанию.")
	flag.Parse()

	// Если указан параметр --make-config, создаем файл конфигурации по умолчанию
	if *makeConfigFlag {
		err := config.Save(configFileName, config.DefaultConfig)
		if err != nil {
			log.Print("Ошибка при создании файла конфигурации:", err)
			return
		}
		log.Print("Создан файл конфигурации по умолчанию ", configFileName)
		return
	}

	// Загрузка конфигурации из файла
	cfg, err := config.Load(configFileName)
	if err != nil {
		log.Print("Ошибка загрузки конфигурации: ", err)
		return
	}
	log.Print("Конфигурация загружена")
	_ = cfg

	// Запускаем камеру
	camera := o2i500.New()
	go camera.Run(config.DefaultConfig.O2i500Addr)

	app := app.New()
	go app.Run()

	// Запускаем графический интерфейс
	gui := guiview.New(&app)
	go gui.Run()

	<-exit
}
