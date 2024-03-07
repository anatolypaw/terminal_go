package entity

// Модель продукта
type Good struct {
	Gtin string `json:"gtin"`
	Desc string `json:"desc"`
	Life int    `json:"life"`
}

// Модель кода
type Code struct {
	Gtin   string `json:"gtin"`
	Serial string `json:"serial"`
	Crypto string `json:"crypto"`
}
