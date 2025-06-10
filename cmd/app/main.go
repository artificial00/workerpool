package main

import (
	"fmt"
	"time"
	"workerpool/pkg/workerpool"
)

// Тут приведен один из тестовых сценариев - просто прогоняем все функции
func main() {
	wp := workerpool.NewWorkerPool()

	time.Sleep(1 * time.Second)

	// тут добавляем нужное количество воркеров
	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()

	// добавляем нужное количество задачек
	for i := 1; i <= 5; i++ {
		wp.SendJob(fmt.Sprintf("задание номер %d", i))
	}

	// удаляем воркера
	wp.RemoveWorker(1)

	// добавляем воркера
	wp.AddWorker()

	// добавляем еще заданий
	for i := 6; i <= 8; i++ {
		wp.SendJob(fmt.Sprintf("задание номер %d", i))
	}

	// опять удаляем воркера
	wp.RemoveWorker(3)

	// симулируем выполнение задачек
	time.Sleep(2 * time.Second)

	// закрываем пул
	wp.Close()
}
