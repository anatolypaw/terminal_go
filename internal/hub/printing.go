package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"terminal/internal/entity"
)

// Возвращает код для печати с сервера маркировки
func (h *Hub) GetCodeForPrint(tname string, gtin string) (entity.Code, error) {
	const op = "hub.GetCodeForPrint"

	type Data struct {
		Tname string
		Gtin  string
	}

	data := Data{
		Tname: tname,
		Gtin:  gtin,
	}

	// Кодирование данных в формат JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Print(op, "ошибка при кодировании JSON:", err)
		return entity.Code{}, err
	}

	// Создание запроса с данными JSON
	url := "http://" + h.HubAddres + "/v1/produce/getCodeForPrint"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(op, "ошибка при создании запроса:", err)
		return entity.Code{}, err
	}

	// Установка заголовков запроса
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	resp, err := h.client.Do(req)
	if err != nil {
		log.Print(op, "ошибка при выполнении запроса:", err)
		return entity.Code{}, err
	}
	defer resp.Body.Close()

	// Чтение содержимого ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("Ошибка при чтении тела ответа:", err)
		return entity.Code{}, err
	}

	var responseData Response
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		log.Print("Ошибка при разборе JSON:", err)
		return entity.Code{}, err
	}

	// Вывод содержимого ответа
	fmt.Println(responseData.Data)

	return entity.Code{}, nil
}
