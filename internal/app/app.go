package app

import (
	"log"
	"terminal/internal/config"
	"terminal/internal/entity"
	"terminal/internal/hub"
	"terminal/internal/o2i500"
	"terminal/internal/savema"
	"time"
)

const (
	MODE_PRODUCE = iota
	MODE_DATACHNGE
	MODE_CANCEL
)

type App struct {
	Cfg *config.Config
	Hub hub.Hub

	Camera  o2i500.O2i500
	Printer savema.Savema

	Mode         int
	Date         int
	SelectedGood entity.Good
}

func New(cfg *config.Config) App {
	return App{
		Cfg: cfg,
	}
}

func (a *App) Run() {
	log.Print("Запуск бизнес логики")
	/* Создаем подключение к хабу */
	a.Hub = hub.New(a.Cfg.P.HubAddr, a.Cfg.P.TerminalName)
	go a.Hub.Run()

	/* Инициализируем используемые устройства */
	if a.Cfg.P.UseCamera {
		log.Print("Используется камера", a.Cfg.P.O2i500Addr)
		a.Camera = o2i500.New(a.Cfg.P.O2i500Addr)
		go a.Camera.Run()
	}

	if a.Cfg.P.UseSavema {
		log.Print("Используется savema: ", a.Cfg.P.SavemaAddr)
		a.Printer = savema.New(a.Cfg.P.SavemaAddr)
		go a.Printer.Run()
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func (a *App) SetMode(m int) {
	a.Mode = m
}

func (a *App) GetMode() int {
	return a.Mode
}

func (a *App) DateDown() {
	a.Date = a.Date - 1
}

func (a *App) DateUp() {
	a.Date = a.Date + 1
}

func (a *App) ModeIsProduce() bool {
	return a.Mode == MODE_PRODUCE
}
