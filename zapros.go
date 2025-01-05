package main

import (
	"fmt"
	"net/http"
)

func zap() {
	// Создаем новый POST-запрос с заголовками
	req, err := http.NewRequest("POST", "http://localhost:8080/ebs/internal/api/createProcessId", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "IntegrationId 52881401")
	req.Header.Set("apiGU", "true")
	req.Header.Set("User-Agent", "")
	req.Header.Set("x-channel", "")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	fmt.Println("Статус ответа:", resp.Status)
	fmt.Println("Тело ответа:", resp.Body)
}
