package workerpool

import (
	"fmt"
	"sync"
)

// WorkerPool структура пула обработчиков
type WorkerPool struct {
	mu        sync.Mutex
	workers   map[int]chan struct{} // мапа воркеров: ключ - айди воркера, значение - канал с пустыми структурами, передающами сигнал о завершении работы воркера
	jobs      chan string           // канал, куда передаются задачи
	wg        sync.WaitGroup
	workerNum int // количество воркеров (служит как счетчик для выдачи айдишников воркерам)
}

// NewWorkerPool функция создания пула обработчиков - сразу инициализируем мапу с воркерами, затем канал с задачами
func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		workers: make(map[int]chan struct{}),
		jobs:    make(chan string),
	}
}

// AddWorker функция добавления воркеров: добавляя воркера, запускаем его горутину, выполняем выданные ему задачи
// и проверяем, что задачи не кончились или что не поступил сигнал о завершении работы воркера
func (wp *WorkerPool) AddWorker() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.workerNum++
	id := wp.workerNum
	stopCh := make(chan struct{})
	wp.workers[id] = stopCh

	wp.wg.Add(1)
	go func(workerID int, stopCh chan struct{}) {
		defer wp.wg.Done()
		for {
			select {
			case job, ok := <-wp.jobs: // проверяем канал на наличие задачек
				if !ok {
					fmt.Printf("Worker c id %d: канал jobs закрыт (тк задачи кончились), выходим\n", workerID)
					return
				}
				fmt.Printf("Worker c id %d делает %s\n", workerID, job)
			case <-stopCh: // пришел сигнал об остановке воркера
				fmt.Printf("Worker c id %d был остановлен\n", workerID)
				return
			}
		}
	}(id, stopCh)

	fmt.Printf("Добавлен worker c id %d\n", id)
	return id
}

// RemoveWorker функция удаления воркера. Удаляем его, закрыв канал stopCh
func (wp *WorkerPool) RemoveWorker(id int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if stopCh, exists := wp.workers[id]; exists {
		close(stopCh)
		delete(wp.workers, id)
		fmt.Printf("Удален worker с id %d\n", id)
	}
}

// SendJob функция добавления задачи (просто посылаем задачу в канал)
func (wp *WorkerPool) SendJob(job string) {
	wp.jobs <- job
}

// Close функция закрытия канала с задачами: закрывает канал, завершая при этом работу оставишхся воркеров.
// wg.Wait добавлен, чтобы  проверить завершение всех горутин перед завершением работы программы
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	fmt.Println("Задачи кончились, закрываем jobs")
	wp.mu.Lock()
	for id, stopCh := range wp.workers {
		close(stopCh)
		delete(wp.workers, id)
	}
	wp.mu.Unlock()
	wp.wg.Wait()
}
