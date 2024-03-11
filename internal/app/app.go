package app

import (
	o2i500 "terminal/internal/O2i500"
	"terminal/internal/config"
	"terminal/internal/entity"
	"terminal/internal/hub"
	"time"
)

const (
	MODE_PRODUCE = iota
	MODE_DATACHNGE
	MODE_CANCEL
)

type App struct {
	Cfg config.Config
	Hub hub.Hub

	Camera o2i500.O2i500

	Mode         int
	Date         int
	SelectedGood entity.Good
}

func New(cfg config.Config) App {
	return App{
		Cfg: cfg,
	}
}

func (a *App) Run() {
	/* Создаем подключение к хабу */
	a.Hub = hub.New(a.Cfg.HubAddr, a.Cfg.TerminalName)
	go a.Hub.Run()

	/* Инициализируем используемые устройства */
	if a.Cfg.UseCamera {
		a.Camera = o2i500.New(a.Cfg.O2i500Addr)
		a.Camera.Run()
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
