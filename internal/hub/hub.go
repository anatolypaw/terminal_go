package hub

import (
	"net/http"
	"time"
)

type Response struct {
	ReqId string `json:"req_id"`
	Ok    bool   `json:"ok"`
	Desc  string `json:"desc"` // Описание результата
	Data  any    `json:"data"`
}

type Hub struct {
	HubAddres string
	client    http.Client

	Connected bool // Флаг, что связь с сервером есть
}

func New(hubAddr string) Hub {
	client := http.Client{
		Timeout: 50 * time.Millisecond,
	}

	return Hub{
		HubAddres: hubAddr,
		client:    client,
	}
}

// Выполняет слежебные операции
func (h *Hub) Run() {
	for {
		// TODO: проверять активность соединения
		time.Sleep(1000 * time.Millisecond)
	}
}
