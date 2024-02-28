package app

import (
	"image/color"
	"time"
)

const (
	MODE_PRODUCE = iota
	MODE_DATACHNGE
	MODE_CANCEL
)

type Good struct {
	Gtin  string
	Desc  string
	Color color.Color
}

type App struct {
	// Режим работы
	// 0 - Производство
	// 1 - Отбраковка
	// 2 - Смена даты
	// 3 - Информация
	Mode int

	Date int

	SelectedGood Good
	Goods        [10]Good
}

func New() App {
	return App{}
}

func (a *App) Run() {
	a.Goods[0] = Good{
		Gtin:  "12313424",
		Desc:  "Творог 330 0%",
		Color: color.RGBA{200, 100, 100, 255},
	}

	a.Goods[1] = Good{
		Gtin:  "12313424",
		Desc:  "Творог 330 5%",
		Color: color.RGBA{255, 255, 100, 255},
	}

	a.Goods[2] = Good{
		Gtin:  "12313424",
		Desc:  "Творог 330 9%",
		Color: color.RGBA{100, 100, 255, 255},
	}

	a.Goods[4] = Good{
		Gtin:  "1231342ad4",
		Desc:  "Йогурт",
		Color: color.RGBA{100, 255, 255, 255},
	}

	a.Goods[8] = Good{
		Gtin:  "12313424eqe",
		Desc:  "Простокваша",
		Color: color.RGBA{100, 255, 255, 255},
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
