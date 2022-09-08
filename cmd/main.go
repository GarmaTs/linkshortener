package main

import (
	chat "github.com/GarmaTs/linkshortener/internal/lesson2/chat"
	//timesender "github.com/GarmaTs/linkshortener/internal/lesson2/time_sender"
)

func main() {
	// Запуск утилыты рассылки даты-времени - решение для первого пункта ДЗ
	// timesender.StartListening(":8001")

	// Запуск чата (обе утилиты одновременно не работают даже на разных портах)
	// Решения для второго и третьего пунктов ДЗ
	chat.StartChatServer(":8000")
}
