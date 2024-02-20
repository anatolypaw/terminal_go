package app

import "time"

type App struct {
	// Режим работы
	// 0 - Производство
	// 1 - Отбраковка
	// 2 - Смена даты
	// 3 - Информация
	Mode int

	I int
}

func New() App {
	return App{}
}

func (a *App) Run() {
	for {
		time.Sleep(1 * time.Second)
	}
}

func (a *App) SetMode(m int) {
	a.I = m
}
