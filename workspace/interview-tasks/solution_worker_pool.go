package tasks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultMaxWorkers = 10
	defaultQueueSize  = 10
	defaultDur        = 250 * time.Millisecond
)

type (
	Task struct{ ID int }

	WorkerPool struct {
		taskQueue  chan Task
		workQueue  chan Task
		shutdownCh chan struct{}
		closeOnce  *sync.Once
		isStopped  *atomic.Bool
		dur        time.Duration
		maxWorkers int
	}

	WorkerPoolOption func(wp *WorkerPool)
)

func WithMaxWorkers(cnt int) WorkerPoolOption {
	if cnt == 0 {
		cnt = defaultMaxWorkers
	}

	return func(wp *WorkerPool) {
		wp.workQueue = make(chan Task, cnt)
	}
}

func WithMaxQueue(n int) WorkerPoolOption {
	if n <= 0 {
		n = defaultQueueSize
	}

	return func(wp *WorkerPool) {
		wp.taskQueue = make(chan Task, n)
	}
}

func WithTimeout(dur time.Duration) WorkerPoolOption {
	if dur <= 0 {
		dur = defaultDur
	}

	return func(wp *WorkerPool) {
		wp.dur = dur
	}
}

func NewWorkerPool(ctx context.Context, opts ...WorkerPoolOption) *WorkerPool {
	wp := &WorkerPool{
		shutdownCh: make(chan struct{}),
		closeOnce:  &sync.Once{},
		isStopped:  &atomic.Bool{},
	}

	for i := range opts {
		opts[i](wp)
	}

	wp.setDefaults()

	// Dispatcher
	go wp.dispatcher(ctx)

	// Workers
	for i := range wp.maxWorkers {
		go wp.worker(ctx, i)
	}

	return wp
}

func (wp *WorkerPool) setDefaults() {
	if wp.taskQueue == nil {
		wp.taskQueue = make(chan Task)
	}

	if wp.workQueue == nil {
		wp.workQueue = make(chan Task)
	}
}

// dispatcher — только он пишет из taskQueue в workQueue
func (wp *WorkerPool) dispatcher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case task, ok := <-wp.taskQueue:
			if !ok {
				return
			}

			if !wp.isStopped.Load() {
				wp.workQueue <- task
			}

		case <-wp.shutdownCh:
			return

		}
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return

		case v, ok := <-wp.workQueue:
			if !ok {
				return
			}

			if wp.isStopped.Load() {
				return
			}

			fmt.Printf("Worker %d handling task %d\n", id, v.ID)
			// Imitating work...
			time.Sleep(2 * time.Second)
			fmt.Printf("Worker %d finished task %d\n", id, v.ID)
		}
	}
}

func (wp *WorkerPool) AddTask(task Task) error {
	return wp.addTask(task)
}

func (wp *WorkerPool) addTask(task Task) error {
	if wp.isStopped.Load() {
		return errors.New("worker pool stopped")
	}

	t := time.NewTicker(wp.dur)
	defer t.Stop()

	select {
	case wp.taskQueue <- task:
		return nil

	case <-t.C:
		return fmt.Errorf("не удалось добавить задачу после %s", wp.dur)
	}

}

// Закрыть пул корректно
func (wp *WorkerPool) Close() {
	wp.closeOnce.Do(func() {
		close(wp.taskQueue)
		close(wp.shutdownCh)
		wp.isStopped.Store(true)
	})
}
