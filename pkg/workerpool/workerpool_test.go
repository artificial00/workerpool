package workerpool

import (
	"fmt"
	"testing"
	"time"
)

// TestWorkerPool_AddAndRemoveWorker тестируем добавление и удаление воркера
func TestWorkerPool_AddAndRemoveWorker(t *testing.T) {
	wp := NewWorkerPool()
	id := wp.AddWorker()
	if id == 0 {
		t.Error("id воркера не должен быть 0")
	}
	wp.RemoveWorker(id)
	wp.mu.Lock()
	_, exists := wp.workers[id]
	wp.mu.Unlock()
	if exists {
		t.Errorf("Воркер с id %d не был удален", id)
	}
}

// TestWorkerPool_CloseAllWorkers тестируем закрытие воркеров по закрытию канала stopCh
func TestWorkerPool_CloseAllWorkers(t *testing.T) {
	wp := NewWorkerPool()
	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()
	stopCh := make(chan struct{})
	go func() {
		wp.Close()
		close(stopCh)
	}()
	select {
	case <-stopCh:
		fmt.Println("Канал закрылся успешно")
	case <-time.After(time.Second):
		t.Error("Close завис, воркеры не были завершены")
	}
}
