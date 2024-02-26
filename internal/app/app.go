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

	SelectedGood *Good
	Goods        [10]Good
}

func New() App {
	return App{}
}

func (a *App) Run() {
	a.Goods[0] = Good{
		Gtin:  "12313424",
		Desc:  "Кефир",
		Color: color.RGBA{100, 255, 100, 255},
	}

	a.Goods[1] = Good{
		Gtin:  "12313424",
		Desc:  "Молоко",
		Color: color.RGBA{255, 255, 100, 255},
	}

	a.Goods[4] = Good{
		Gtin:  "12313424",
		Desc:  "Йогурт",
		Color: color.RGBA{100, 255, 255, 255},
	}

	for {
		a.SelectedGood = &Good{
			Gtin:  "1231241412",
			Desc:  "Молоко",
			Color: color.RGBA{100, 100, 200, 255},
		}

		time.Sleep(1 * time.Second)

		a.SelectedGood = nil

		time.Sleep(1 * time.Second)
	}
}

func (a *App) SetMode(m int) {
	a.Mode = m
}

func (a *App) GetMode() int {
	return a.Mode
}
